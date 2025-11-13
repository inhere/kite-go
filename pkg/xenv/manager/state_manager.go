package manager

import (
	"fmt"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/goccy/go-json"
	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/jsonutil"
	"github.com/gookit/goutil/strutil"
	"github.com/inhere/kite-go/pkg/xenv/models"
	"github.com/inhere/kite-go/pkg/xenv/xenvcom"
)

// StateManager manages the state data of the environment
type StateManager struct {
	init bool
	// batch mode
	batchMode bool
	// all merged state data.
	//  - 数据按 global, direnv, session 的顺序合并进来
	merged *models.ActivityState
	// session current session state data. NOTE: 只有在 HOOK SHELL 中才会生效
	session *models.ActivityState
	// global state data from global state file xenvcom.GlobalStateFile
	global *models.ActivityState
	// directory states. 从当前目录开始，会向上级目录递归查找 .xenv.toml
	//  - index 越大的优先级越高，距离当前目录越近
	//  - 当前只会查找最近的一个 .xenv.toml 文件 TODO 后续支持多个
	dirStates  []*models.ActivityState
	envrcFiles []string // TODO
}

// NewStateManager creates a new StateManager
func NewStateManager() *StateManager {
	globalFile := fsutil.ExpandHome(xenvcom.GlobalStateFile)
	sessionFile := fsutil.ExpandHome(xenvcom.SessionFile())

	return &StateManager{
		merged: models.NewActivityState("MERGED"),
		global:  models.NewActivityState(globalFile),
		session: models.NewActivityState(sessionFile),
		dirStates: make([]*models.ActivityState, 0),
	}
}

// Init initializes the state manager
func (m *StateManager) Init() error {
	if m.init {
		return nil
	}
	m.init = true
	return m.LoadStateFiles()
}

// GlobalFile returns the global state file path
func (m *StateManager) GlobalFile() string { return m.global.File }

// SetBatchMode sets the batch mode flag
func (m *StateManager) SetBatchMode(enabled bool) { m.batchMode = enabled }

// endregion
// region SDK state management
//

// UseSDKsWithParams activates multiple SDKs with params
func (m *StateManager) UseSDKsWithParams(ps *models.ActivateSDKsParams) error {
	if err := m.requireInit(); err != nil {
		return err
	}

	// update 合并数据
	m.merged.AddSDKs(ps.AddSdks).AddEnvs(ps.AddEnvs).DelPaths(ps.RemPaths).AddPaths(ps.AddPaths)

	switch ps.OpFlag {
	case models.OpFlagGlobal:
		m.global.AddSDKs(ps.AddSdks)
	case models.OpFlagDirenv:
		ds := m.DirenvOrNew()
		ds.AddSDKs(ps.AddSdks)
		// save to session dirStates
		m.session.AddDirState(ds)
	default:
		m.session.AddSDKs(ps.AddSdks)
	}

	// Save the state file
	if !m.batchMode {
		return m.SaveStateFile()
	}
	return nil
}

// ActivateSDK activates a specific SDK version
func (m *StateManager) ActivateSDK(name, version string, opFlag models.OpFlag) error {
	if err := m.requireInit(); err != nil {
		return err
	}

	// update 合并数据
	m.merged.SDKs[name] = version

	switch opFlag {
	case models.OpFlagGlobal:
		m.global.SDKs[name] = version
	case models.OpFlagDirenv:
		m.DirenvOrNew().SDKs[name] = version
	default:
		m.session.SDKs[name] = version
	}

	// Save the state file
	if !m.batchMode {
		return m.SaveStateFile()
	}
	return nil
}

// DelSDKsWithEnvsPaths deletes multiple tools and with envs, paths
func (m *StateManager) DelSDKsWithEnvsPaths(names, envs, paths []string, opFlag models.OpFlag) error {
	if err := m.requireInit(); err != nil {
		return err
	}

	// update 合并数据
	m.merged.DelSDKsEnvsPaths(names, envs, paths)

	switch opFlag {
	case models.OpFlagGlobal:
		m.global.DelSDKs(names)
	case models.OpFlagDirenv:
		if ds := m.Nearest(); ds != nil {
			ds.DelSDKs(names)
		}
	default:
		m.session.DelSDKs(names)
	}

	// Save the state file
	if !m.batchMode {
		return m.SaveStateFile()
	}
	return nil
}

// endregion
// region Env state management
//

// AddEnvs sets multiple environment variables
func (m *StateManager) AddEnvs(envs map[string]string, opFlag models.OpFlag) error {
	if err := m.requireInit(); err != nil {
		return err
	}

	// update 合并数据
	m.merged.AddEnvs(envs)

	switch opFlag {
	case models.OpFlagGlobal:
		m.global.AddEnvs(envs)
	case models.OpFlagDirenv:
		m.DirenvOrNew().AddEnvs(envs)
	default:
		m.session.AddEnvs(envs)
	}

	// Save the state file
	if !m.batchMode {
		return m.SaveStateFile()
	}
	return nil
}

// SetEnv sets an environment variable
func (m *StateManager) SetEnv(name, value string, opFlag models.OpFlag) error {
	if err := m.requireInit(); err != nil {
		return err
	}

	m.merged.Envs[name] = value

	switch opFlag {
	case models.OpFlagGlobal:
		m.global.Envs[name] = value
	case models.OpFlagDirenv:
		m.DirenvOrNew().Envs[name] = value
	default:
		m.session.Envs[name] = value
	}

	if !m.batchMode {
		return m.SaveStateFile()
	}
	return nil
}

// UnsetEnv unsets an environment variable
func (m *StateManager) UnsetEnv(name string, opFlag models.OpFlag) error {
	if err := m.requireInit(); err != nil {
		return err
	}

	// update 删除合并数据
	// check exists
	if _, exists := m.merged.Envs[name]; exists {
		delete(m.merged.Envs, name)
	} else {
		return fmt.Errorf("environment variable %s is not exists", name)
	}

	switch opFlag {
	case models.OpFlagGlobal:
		delete(m.global.Envs, name)
	case models.OpFlagDirenv:
		if ds := m.Nearest(); ds != nil {
			delete(ds.Envs, name)
		}
	default:
		delete(m.session.Envs, name)
	}

	if !m.batchMode {
		return m.SaveStateFile()
	}
	return nil
}

// endregion
// region PATH state management
//

// AddPaths adds multiple paths to the PATH environment variable
func (m *StateManager) AddPaths(paths []string, opFlag models.OpFlag) (err error) {
	if err = m.requireInit(); err != nil {
		return err
	}

	m.merged.AddPaths(paths)

	switch opFlag {
	case models.OpFlagGlobal:
		m.global.AddPaths(paths)
	case models.OpFlagDirenv:
		m.DirenvOrNew().AddPaths(paths)
	default:
		m.session.AddPaths(paths)
	}

	if !m.batchMode {
		return m.SaveStateFile()
	}
	return
}

// AddPath adds a path to the PATH environment variable
func (m *StateManager) AddPath(path string, opFlag models.OpFlag) error {
	if err := m.requireInit(); err != nil {
		return err
	}

	// update 添加合并数据
	m.merged.AddPath(path)

	switch opFlag {
	case models.OpFlagGlobal:
		m.global.AddPath(path)
	case models.OpFlagDirenv:
		m.DirenvOrNew().AddPath(path)
	default:
		m.session.AddPath(path)
	}

	if !m.batchMode {
		return m.SaveStateFile()
	}
	return nil
}

// DelPath removes a path from the PATH environment variable
func (m *StateManager) DelPath(path string, opFlag models.OpFlag) error {
	if err := m.requireInit(); err != nil {
		return err
	}

	m.merged.DelPath(path)

	switch opFlag {
	case models.OpFlagGlobal:
		m.global.DelPath(path)
	case models.OpFlagDirenv:
		if ds := m.Nearest(); ds != nil {
			ds.DelPath(path)
		}
	default:
		m.session.DelPath(path)
	}

	if !m.batchMode {
		return m.SaveStateFile()
	}
	return nil
}

// endregion
// region Load/Save State Files
//

// LoadStateFiles loads the global and dir state from file
func (m *StateManager) LoadStateFiles() (err error) {
	// Load the global state
	if err = m.loadTomlStateFile(m.global); err != nil {
		return fmt.Errorf("failed to load global state: %w", err)
	}

	// Merge the global state into the session state
	m.merged.Merge(m.global)

	// Load the direnv state
	if err = m.LoadDirEnvState(); err != nil {
		return fmt.Errorf("failed to load direnv state: %w", err)
	}

	// Load the session state
	if xenvcom.InHookShell() {
		if fsutil.IsFile(m.session.File) {
			if err = m.loadJsonStateFile(m.session); err == nil {
				m.merged.Merge(m.session)
			}
		}
		m.session.Shell = xenvcom.HookShell()
	}

	return
}

func (m *StateManager) LoadDirEnvState() error {
	// Get the current working directory
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %w", err)
	}

	// Check for .xenv.toml file in the current directory and parent directories up to the root
	xenvTomlPath := fsutil.FindOneInParentDirs(wd, xenvcom.LocalStateFile)
	if xenvTomlPath != "" {
		xenvcom.Debugf("Found .xenv.toml at: %s\n", xenvTomlPath)
		// Process the .xenv.toml file
		if err1 := m.processDirenvToml(xenvTomlPath); err1 != nil {
			return fmt.Errorf("failed to process .xenv.toml: %w", err1)
		}
	}

	// Check for .envrc file in the current directory and parent directories up to the root
	fileName := strutil.OrCond(xenvcom.IsHookBash(), ".envrc", ".envrc.ps1")
	envrcPath := fsutil.FindOneInParentDirs(wd, fileName)
	if envrcPath != "" {
		xenvcom.Debugf("Found envrc at: %s\n", envrcPath)
		m.envrcFiles = append(m.envrcFiles, envrcPath)
	}

	return nil
}

// processDirenvToml processes an .xenv.toml file
func (m *StateManager) processDirenvToml(filePath string) error {
	dirState := models.NewActivityState(filePath)
	err := m.loadTomlStateFile(dirState)
	if err != nil {
		return fmt.Errorf("failed to load dir state: %w", err)
	}

	m.merged.Merge(dirState)
	m.dirStates = append(m.dirStates, dirState)
	return nil
}

// loads the xenv state from TOML file
func (m *StateManager) loadTomlStateFile(ptr *models.ActivityState) error {
	// Check if file exists
	if _, err := os.Stat(ptr.File); os.IsNotExist(err) {
		return nil
	}

	// Read the TOML file
	data, err := os.ReadFile(ptr.File)
	if err != nil {
		return err
	}
	return toml.Unmarshal(data, ptr)
}

// loads the xenv state from TOML file
func (m *StateManager) loadJsonStateFile(ptr *models.ActivityState) error {
	// Check if file exists
	if _, err := os.Stat(ptr.File); os.IsNotExist(err) {
		return nil
	}
	xenvcom.Debugf("Loading session file: %s\n", ptr.File)

	// Read the JSON file
	data, err := os.ReadFile(ptr.File)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, ptr)
}

// SaveStateFile saves the global state to file
func (m *StateManager) SaveStateFile() error {
	if m.global.HasUpdate {
		if err := m.saveStateFile(m.global); err != nil {
			return err
		}
		m.global.HasUpdate = false
	}

	// direnv states
	for _, state := range m.dirStates {
		if state.HasUpdate {
			if err := m.saveStateFile(state); err != nil {
				return err
			}
			state.HasUpdate = false
		}
	}

	// 会话数据: 只有在 HOOK SHELL 中才会生效
	if xenvcom.InHookShell() && m.session.HasUpdate {
		m.session.HasUpdate = false
		return m.saveStateFile(m.session)
	}
	return nil
}

func (m *StateManager) saveStateFile(state *models.ActivityState) error {
	xenvcom.Debugf("Saving state file: %s\n", state.File)

	// is session state file. save as json
	if state.IsSession() {
		if err := fsutil.MkParentDir(state.File); err != nil {
			return err
		}

		state.UpdatedAt = time.Now()
		if state.CreatedAt.IsZero() {
			state.CreatedAt = state.UpdatedAt
		}
		return jsonutil.WritePretty(state.File, state)
	}
	return NewTomlUpdater().Update(state)
}

// endregion
// region Helper methods
//

func (m *StateManager) requireInit() error {
	if !m.init {
		return fmt.Errorf("please call Init() before using the state manager")
	}
	return nil
}

func (m *StateManager) Merged() *models.ActivityState { return m.merged }

// Global returns the global activity state
func (m *StateManager) Global() *models.ActivityState { return m.global }

// DirStates returns the direnv activity states
func (m *StateManager) DirStates() []*models.ActivityState { return m.dirStates }

// Session returns the session activity state
func (m *StateManager) Session() *models.ActivityState { return m.session }

// Nearest returns the nearest direnv activity state
func (m *StateManager) Nearest() *models.ActivityState {
	if len(m.dirStates) > 0 {
		return m.dirStates[len(m.dirStates)-1]
	}
	return nil
}

// DirenvOrNew returns the nearest direnv activity state or a new one at workdir
func (m *StateManager) DirenvOrNew() *models.ActivityState {
	if len(m.dirStates) > 0 {
		return m.dirStates[len(m.dirStates)-1]
	}

	// TODO 输出提示，确认是否创建 .xenv.toml 文件
	de := models.NewActivityState(xenvcom.LocalStateFile)
	m.dirStates = append(m.dirStates, de)
	return de
}
