package service

import (
	"fmt"
	"os"
	"strings"

	"github.com/gookit/goutil/maputil"
	"github.com/gookit/goutil/strutil"
	"github.com/gookit/goutil/x/ccolor"
	"github.com/inhere/kite-go/pkg/util"
	"github.com/inhere/kite-go/pkg/xenv/manager"
	"github.com/inhere/kite-go/pkg/xenv/models"
	"github.com/inhere/kite-go/pkg/xenv/shell"
)

// EnvService handles environment variable and PATH management
type EnvService struct {
	config *models.Configuration
	state  *manager.StateManager
	// envMgr *manager.EnvManager
}

// NewEnvService creates a new EnvService
func NewEnvService(config *models.Configuration, state *manager.StateManager) *EnvService {
	return &EnvService{
		config: config,
		state:  state,
		// envMgr: manager.NewEnvManager(),
	}
}

// IsSessionEnv 判断当前是否在shell hook环境
func (s *EnvService) IsSessionEnv() bool {
	return util.InHookShell()
}

func (s *EnvService) GlobalState() *models.ActivityState {
	return s.state.Global()
}

func (s *EnvService) SessionState() *models.ActivityState {
	return s.state.Session()
}

// endregion
// region Shell Hook Init
//

// WriteHookToProfile installs the hook script to the user's profile
func (s *EnvService) WriteHookToProfile(st shell.ShellType, pwshProfile string) error {
	gen := shell.NewScriptGenerator(st, s.config)
	if util.InHookShell() {
		ccolor.Infoln("The hook script is already installed in the current shell")
		return nil
	}

	return gen.InstallToProfile(pwshProfile)
}

// GenHookScripts generates Shell hook init scripts
func (s *EnvService) GenHookScripts(st shell.ShellType) (string, error) {
	gen := shell.NewScriptGenerator(st, s.config)

	return gen.GenHookScripts(s.state.Global())
}

// endregion
// region ENV management
//

// SetEnv sets an environment variable
func (s *EnvService) SetEnv(name, value string, global bool) (script string, err error) {
	// Generate shell eval scripts
	gen, err1 := getShellGenerator(s.config)
	if err1 != nil {
		return "", err1
	}

	name = strings.ToUpper(name)
	if !strutil.IsVarName(name) {
		return "", fmt.Errorf("invalid environment variable name: %s", name)
	}

	// 在shell hook环境中, 生成 ENV set 脚本
	if gen != nil {
		script = gen.GenSetEnv(name, value)
	} else {
		ccolor.Warnln("TIP: The operation will not take effect, please setup the SHELL HOOK first.")
	}

	// TIP: 设置程序内部 ENV 没有意义
	// if err := os.Setenv(name, value); err != nil {
	// 	return "", fmt.Errorf("failed to set environment variable: %w", err)
	// }

	// Add to activity state data
	err = s.state.SetEnv(name, value, global)
	return
}

// UnsetEnvs unsets multi environment variables
func (s *EnvService) UnsetEnvs(names []string, global bool) (script string, err error) {
	var sb strings.Builder
	// Generate shell eval scripts
	gen, err1 := getShellGenerator(s.config)
	if err1 != nil {
		return "", err1
	}
	if gen == nil {
		ccolor.Warnln("TIP: The operation will not take effect, please setup the SHELL HOOK first.")
	}

	s.state.SetBatchMode(true)
	defer s.state.SetBatchMode(false)
	for _, name := range names {
		name = strings.ToUpper(name)
		if val := os.Getenv(name); val == "" {
			ccolor.Warnf("ENV var not found: %s\n", name)
		}

		// 在shell hook环境中, 生成ENV set脚本
		if gen != nil {
			sb.WriteString(gen.GenUnsetEnv(name))
		}

		err = s.state.UnsetEnv(name, global)
		if err != nil {
			return "", err
		}
	}

	if global {
		err = s.state.SaveStateFile()
	}
	return sb.String(), err
}

// GlobalEnv lists environment variables in the global scope
func (s *EnvService) GlobalEnv() map[string]string {
	// Return the global environment variables
	return maputil.MergeStrMap(s.config.GlobalEnv, s.state.Global().Envs)
}

// SessionEnv lists environment variables in the current session
func (s *EnvService) SessionEnv() map[string]string {
	return s.state.Session().Envs
}

// endregion
// region PATH management
//

// AddPath adds a path to the PATH environment variable
func (s *EnvService) AddPath(path string, global bool) (script string, err error) {
	normalizedPath := util.NormalizePath(path)

	// Check if path exists
	if _, err := os.Stat(normalizedPath); os.IsNotExist(err) {
		return "", fmt.Errorf("path does not exist: %s", normalizedPath)
	}

	// Add to session PATH
	currentPath := os.Getenv("PATH")
	pathList := util.SplitPath(currentPath)

	// Check if path already exists
	for _, p := range pathList {
		if p == normalizedPath {
			return "", fmt.Errorf("path already exists in PATH: %s", normalizedPath)
		}
	}

	// Generate shell eval scripts
	gen, err1 := getShellGenerator(s.config)
	if err1 != nil {
		return "", err1
	}

	// Add the path to the beginning of PATH (highest priority)
	// newPathList := append([]string{normalizedPath}, pathList...)

	// 在shell hook环境中, 生成 ENV set 脚本
	if gen != nil {
		script = gen.GenAddPath(normalizedPath)
	} else {
		ccolor.Warnln("TIP: The operation will not take effect, please setup the SHELL HOOK first.")
	}

	// Add to activity state
	err = s.state.AddPath(normalizedPath, global)
	return
}

// RemovePath removes a path from the PATH environment variable
func (s *EnvService) RemovePath(path string, global bool) (script string, err error) {
	// Normalize the path
	normalizedPath := util.NormalizePath(path)
	pathList := util.SplitPath(os.Getenv("PATH"))

	found := false
	var newPaths []string

	// Remove from session PATH
	for _, p := range pathList {
		if p != normalizedPath {
			newPaths = append(newPaths, p)
		} else {
			found = true
		}
	}
	if !found {
		ccolor.Warnf("path not found in PATH: %s", normalizedPath)
		// Remove from activity state
		err = s.state.RemovePath(normalizedPath, global)
		return "", err
	}

	// Generate shell eval scripts
	gen, err1 := getShellGenerator(s.config)
	if err1 != nil {
		return "", err1
	}

	// 在shell hook环境中, 生成 ENV set 脚本
	if gen != nil {
		script = gen.GenSetPath(newPaths)
	} else {
		ccolor.Warnln("TIP: The operation will not take effect, please setup the SHELL HOOK first.")
	}

	// Remove from activity state
	err = s.state.RemovePath(normalizedPath, global)
	return
}

// ListPaths lists PATH entries
func (s *EnvService) ListPaths() []models.PathEntry {
	var paths []models.PathEntry
	for _, entry := range s.state.Global().Paths {
		paths = append(paths, models.PathEntry{
			Path:     entry,
			Priority: 0,
			IsActive: true,
			Scope:    "global",
		})
	}

	for _, entry := range s.state.Session().Paths {
		paths = append(paths, models.PathEntry{
			Path:     entry,
			Priority: 0,
			IsActive: true,
			Scope:    "session",
		})
	}
	return paths
}

// SearchPath searches for a path in PATH
func (s *EnvService) SearchPath(path string) []string {
	normalizedPath := util.NormalizePath(path)
	var matches []string

	// Search in active paths
	for _, p := range s.state.Global().Paths {
		if strings.Contains(p, normalizedPath) {
			matches = append(matches, p)
		}
	}

	// Also search in current system PATH
	currentPath := os.Getenv("PATH")
	pathList := util.SplitPath(currentPath)
	for _, p := range pathList {
		if strings.Contains(p, normalizedPath) {
			matches = append(matches, p)
		}
	}

	return matches
}
