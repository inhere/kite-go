package manager

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/gookit/goutil"
	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/strutil"
	"github.com/inhere/kite-go/pkg/util"
	"github.com/inhere/kite-go/pkg/xenv/models"
)

// StateManager manages the state data of the environment
type StateManager struct {
	init      bool
	batchMode bool
	// global state file path
	stateFile string
	// session merged state data
	session *models.ActivityState
	// global state data from global state file models.GlobalStateFile
	global *models.ActivityState
	// directory states. 从当前目录开始，会向上级目录递归查找 .xenv.toml
	//  - 当前只会查找最近的一个 .xenv.toml 文件 TODO 后续支持多个
	dirStates  []*models.ActivityState
	envrcFiles []string
}

// NewStateManager creates a new StateManager
func NewStateManager() *StateManager {
	return &StateManager{
		stateFile: fsutil.ExpandHome(models.GlobalStateFile),
		session:   models.NewActivityState("SESSION"),
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

// StateFile returns the state file path
func (m *StateManager) StateFile() string { return m.stateFile }

// SetBatchMode sets the batch mode flag
func (m *StateManager) SetBatchMode(enabled bool) { m.batchMode = enabled }

// endregion
// region Tool state management
//

// UseToolsWithParams activates multiple tools with params
func (m *StateManager) UseToolsWithParams(ps *models.ActivateToolsParams, global bool) error {
	if err := m.Init(); err != nil {
		return err
	}

	if global {
		m.global.RemovePaths(ps.RemPaths)
		m.global.AddToolsWithEnvsPaths(ps.AddTools, ps.AddEnvs, ps.AddPaths)
		// Save the global state
		if !m.batchMode {
			return m.SaveStateFile()
		}
	} else {
		m.session.RemovePaths(ps.RemPaths)
		m.session.AddToolsWithEnvsPaths(ps.AddTools, ps.AddEnvs, ps.AddPaths)
	}

	return nil
}

// UseToolsWithEnvsPaths activates multiple tools and with envs, paths
func (m *StateManager) UseToolsWithEnvsPaths(tools, envs map[string]string, paths []string, global bool) error {
	if err := m.Init(); err != nil {
		return err
	}

	if global {
		m.global.AddToolsWithEnvsPaths(tools, envs, paths)
		// Save the global state
		if !m.batchMode {
			return m.SaveStateFile()
		}
	} else {
		m.session.AddToolsWithEnvsPaths(tools, envs, paths)
	}

	return nil
}

// ActivateTool activates a specific tool version
func (m *StateManager) ActivateTool(name, version string, global bool) error {
	if err := m.Init(); err != nil {
		return err
	}

	if global {
		m.global.SDKs[name] = version
		// Save the global state
		if !m.batchMode {
			return m.SaveStateFile()
		}
	} else {
		m.session.SDKs[name] = version
	}
	return nil
}

// DelToolsWithEnvsPaths deletes multiple tools and with envs, paths
func (m *StateManager) DelToolsWithEnvsPaths(tools, envs, paths []string, global bool) error {
	if err := m.Init(); err != nil {
		return err
	}

	if global {
		m.global.DelToolsWithEnvsPaths(tools, envs, paths)
		// Save the global state
		if !m.batchMode {
			return m.SaveStateFile()
		}
	} else {
		m.session.DelToolsWithEnvsPaths(tools, envs, paths)
	}
	return nil
}

// DeactivateTool deactivates a specific tool version
func (m *StateManager) DeactivateTool(name string, global bool) error {
	if err := m.Init(); err != nil {
		return err
	}

	if global {
		if m.global.RemoveTool(name) {
			// Save the global state
			if !m.batchMode {
				return m.SaveStateFile()
			}
		}
		return fmt.Errorf("tool %q was never activated in global", name)
	}

	if m.session.RemoveTool(name) {
		return nil
	}
	return fmt.Errorf("tool %q was never activated in session", name)
}

// endregion
// region Env state management
//

// AddEnvs sets multiple environment variables
func (m *StateManager) AddEnvs(envs map[string]string, global bool) error {
	if err := m.Init(); err != nil {
		return err
	}

	if global {
		for name, value := range envs {
			m.global.Envs[name] = value
		}
		// Save the global state
		if !m.batchMode {
			return m.SaveStateFile()
		}
		return nil
	}

	for name, value := range envs {
		m.session.Envs[name] = value
	}
	return nil
}

// SetEnv sets an environment variable
func (m *StateManager) SetEnv(name, value string, global bool) error {
	if err := m.Init(); err != nil {
		return err
	}

	// Set global env
	if global {
		m.global.Envs[name] = value
		// Save the global state
		if !m.batchMode {
			return m.SaveStateFile()
		}
	} else {
		// Set session env
		m.session.Envs[name] = value
	}
	return nil
}

// UnsetEnv unsets an environment variable
func (m *StateManager) UnsetEnv(name string, global bool) error {
	if err := m.Init(); err != nil {
		return err
	}

	// Unset global env
	if global {
		// check exists
		if _, exists := m.global.Envs[name]; exists {
			// return fmt.Errorf("environment variable %s is not currently set", name)
			delete(m.global.Envs, name)
			// Save the global state
			if !m.batchMode {
				return m.SaveStateFile()
			}
		}
	} else {
		// check exists
		if _, exists := m.session.Envs[name]; exists {
			// return fmt.Errorf("environment variable %s is not currently set", name)
			delete(m.session.Envs, name)
		}
	}
	return nil
}

// endregion
// region PATH state management
//

// AddPaths adds multiple paths to the PATH environment variable
func (m *StateManager) AddPaths(paths []string, global bool) (err error) {
	if err = m.Init(); err != nil {
		return err
	}

	if global {
		for _, path := range paths {
			m.global.AddActivePath(path)
		}
		// Save the global state
		if !m.batchMode {
			return m.SaveStateFile()
		}
		return
	}

	for _, path := range paths {
		m.session.AddActivePath(path)
	}
	return
}

// AddPath adds a path to the PATH environment variable
func (m *StateManager) AddPath(path string, global bool) error {
	if err := m.Init(); err != nil {
		return err
	}

	// Add to global
	if global {
		m.global.AddActivePath(path)
		// Save the global state
		if !m.batchMode {
			return m.SaveStateFile()
		}
		return nil
	}

	// Add to session
	m.session.AddActivePath(path)
	return nil
}

// RemovePath removes a path from the PATH environment variable
func (m *StateManager) RemovePath(path string, global bool) error {
	if err := m.Init(); err != nil {
		return err
	}

	if global {
		newPaths := make([]string, 0, len(m.global.Paths))
		for _, p := range m.global.Paths {
			if p != path {
				newPaths = append(newPaths, p)
			}
		}
		m.global.Paths = newPaths
		// Save the global state
		if !m.batchMode {
			return m.SaveStateFile()
		}
	}

	newPaths := make([]string, 0, len(m.session.Paths))
	for _, p := range m.session.Paths {
		if p != path {
			newPaths = append(newPaths, p)
		}
	}
	m.session.Paths = newPaths
	return nil
}

// endregion
// region Load/Save State Files
//

// LoadStateFiles loads the global and dir state from file
func (m *StateManager) LoadStateFiles() error {
	// Load the global state
	m.global = models.NewActivityState(m.stateFile)
	err := m.loadStateFile(m.stateFile, m.global)
	if err != nil {
		return fmt.Errorf("failed to load global state: %w", err)
	}

	// Merge the global state into the session state
	m.session.Merge(m.global)

	// Load the direnv state
	err = m.LoadDirEnvState()
	if err != nil {
		return fmt.Errorf("failed to load direnv state: %w", err)
	}

	return nil
}

func (m *StateManager) LoadDirEnvState() error {
	// Get the current working directory
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %w", err)
	}

	// Check for .xenv.toml file in the current directory and parent directories up to the root
	xenvTomlPath := fsutil.FindOneInParentDirs(wd, ".xenv.toml")
	if xenvTomlPath != "" {
		fmt.Printf("Found .xenv.toml at: %s\n", xenvTomlPath)
		// Process the .xenv.toml file
		if err1 := m.processDirenvToml(xenvTomlPath); err1 != nil {
			return fmt.Errorf("failed to process .xenv.toml: %w", err1)
		}
	}

	// Check for .envrc file in the current directory and parent directories up to the root
	fileName := strutil.OrCond(util.IsHookBash(), ".envrc", ".envrc.ps1")
	envrcPath := fsutil.FindOneInParentDirs(wd, fileName)
	if envrcPath != "" {
		fmt.Printf("Found envrc at: %s\n", envrcPath)
		m.envrcFiles = append(m.envrcFiles, envrcPath)
	}

	return nil
}

// processXenvToml processes an .xenv.toml file
func (m *StateManager) processDirenvToml(filePath string) error {
	dirState := models.NewActivityState(filePath)
	err := m.loadStateFile(filePath, dirState)
	if err != nil {
		return fmt.Errorf("failed to load dir state: %w", err)
	}

	m.session.Merge(dirState)
	m.dirStates = append(m.dirStates, dirState)
	return nil
}

// loads the xenv state from TOML file
func (m *StateManager) loadStateFile(stateFile string, ptr *models.ActivityState) error {
	// Check if file exists
	if _, err := os.Stat(stateFile); os.IsNotExist(err) {
		return nil
	}

	// Read the file
	data, err := os.ReadFile(stateFile)
	if err != nil {
		return err
	}

	// Unmarshal the TOML
	return toml.Unmarshal(data, ptr)
}

// UpdateStateFile edits a TOML file
func (m *StateManager) UpdateStateFile(state *models.ActivityState) error {
	if state.ChangeData == nil {
		return nil
	}

	return nil
}

// SaveStateFile saves the global state to file
func (m *StateManager) SaveStateFile() error {
	if err := fsutil.MkParentDir(m.stateFile); err != nil {
		return err
	}

	// Update timestamps
	// m.global.UpdatedAt = time.Now()
	// if m.global.CreatedAt.IsZero() {
	// 	m.global.CreatedAt = m.global.UpdatedAt
	// }

	// Marshal to JSON
	data, err := json.MarshalIndent(m.global, "", "  ")
	if err != nil {
		return err
	}

	// Write to file
	return os.WriteFile(m.stateFile, data, 0644)
}

// endregion
// region Helper methods
//

func (m *StateManager) requireInit() {
	goutil.PanicErr(m.Init())
}

// Global returns the global activity state
func (m *StateManager) Global() *models.ActivityState {
	return m.global
}

// Session returns the session activity state
func (m *StateManager) Session() *models.ActivityState {
	return m.session
}

// GetActiveTool returns the active tool version for a given name
func (m *StateManager) GetActiveTool(name string, global bool) string {
	if global {
		return m.global.SDKs[name]
	}
	return m.session.SDKs[name]
}
