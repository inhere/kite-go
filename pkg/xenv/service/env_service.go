package service

import (
	"fmt"
	"os"
	"strings"

	"github.com/gookit/goutil/maputil"
	"github.com/gookit/goutil/x/ccolor"
	"github.com/inhere/kite-go/internal/util"
	"github.com/inhere/kite-go/pkg/xenv/manager"
	"github.com/inhere/kite-go/pkg/xenv/models"
	"github.com/inhere/kite-go/pkg/xenv/shell"
)

// EnvService handles environment variable and PATH management
type EnvService struct {
	config *models.Configuration
	state  *manager.StateManager
}

// NewEnvService creates a new EnvService
func NewEnvService(config *models.Configuration, state *manager.StateManager) *EnvService {
	return &EnvService{
		config: config,
		state:  state,
	}
}

// IsSessionEnv 判断当前是否在shell hook环境
func (m *EnvService) IsSessionEnv() bool {
	return shell.InHookShell()
}

// getShellGenerator 获取当前shell的脚本生成器. 注意：不在shell hook环境，会返回nil
func (m *EnvService) getShellGenerator() (*shell.XenvScriptGenerator, error) {
	// hookShell 不为空表明在shell hook环境中
	hookShell := shell.HookShell()
	if hookShell == "" {
		return nil, nil
	}

	shellType, err := shell.TypeFromString(hookShell)
	if err != nil {
		return nil, err
	}

	return shell.NewScriptGenerator(shellType, m.config), nil
}

// SetEnv sets an environment variable
func (m *EnvService) SetEnv(name, value string, global bool) (script string, err error) {
	// Generate shell eval scripts
	gen, err1 := m.getShellGenerator()
	if err1 != nil {
		return "", err1
	}
	// 在shell hook环境中, 生成 ENV set 脚本
	if gen != nil {
		script = gen.GenEnvSet(name, value)
	} else {
		ccolor.Magentaln("TIP: NOT IN SHELL HOOK. Will not take effect in current shell")
	}

	// TIP: 设置程序内部 ENV 没有意义
	// if err := os.Setenv(name, value); err != nil {
	// 	return "", fmt.Errorf("failed to set environment variable: %w", err)
	// }

	// Add to activity state data
	err = m.state.SetEnv(name, value, global)
	return
}

// endregion
// region ENV management
//

// UnsetEnv unsets an environment variable
func (m *EnvService) UnsetEnv(name string, global bool) error {
	// Remove from session
	if err := os.Unsetenv(name); err != nil {
		return fmt.Errorf("failed to unset environment variable: %w", err)
	}

	// Remove from activity state
	return m.state.UnsetEnv(name, global)
}

// UnsetEnvs unsets multi environment variables
func (m *EnvService) UnsetEnvs(names []string, global bool) (script string, err error) {
	var sb strings.Builder
	// Generate shell eval scripts
	gen, err1 := m.getShellGenerator()
	if err1 != nil {
		return "", err1
	}
	if gen == nil {
		ccolor.Magentaln("TIP: NOT IN SHELL HOOK. Will not take effect in current shell")
	}

	m.state.SetBatchMode(true)
	defer m.state.SetBatchMode(false)
	for _, name := range names {
		// 在shell hook环境中, 生成ENV set脚本
		if gen != nil {
			sb.WriteString(gen.GenEnvUnset(name))
		}
		err = m.state.UnsetEnv(name, global)
		if err != nil {
			return "", err
		}
	}

	if global {
		err = m.state.SaveGlobalState()
	}
	return sb.String(), err
}

// GlobalEnv lists environment variables in the global scope
func (m *EnvService) GlobalEnv() map[string]string {
	// Return the global environment variables
	return maputil.MergeStrMap(m.config.GlobalEnv, m.state.Global().ActiveEnv)
}

// SessionEnv lists environment variables in the current session
func (m *EnvService) SessionEnv() map[string]string {
	return m.state.Session().ActiveEnv
}

// endregion
// region PATH management
//

// AddPath adds a path to the PATH environment variable
func (m *EnvService) AddPath(path string, global bool) error {
	// Normalize the path
	normalizedPath := util.NormalizePath(path)

	// Check if path exists
	if _, err := os.Stat(normalizedPath); os.IsNotExist(err) {
		return fmt.Errorf("path does not exist: %s", normalizedPath)
	}

	// Add to session PATH
	currentPath := os.Getenv("PATH")
	pathList := util.SplitPathList(currentPath)

	// Check if path already exists
	for _, p := range pathList {
		if p == normalizedPath {
			return fmt.Errorf("path already exists in PATH: %s", normalizedPath)
		}
	}

	// Add the path to the beginning of PATH (highest priority)
	newPathList := append([]string{normalizedPath}, pathList...)
	newPath := util.JoinPathList(newPathList)

	// Set the new PATH environment variable
	if err := os.Setenv("PATH", newPath); err != nil {
		return fmt.Errorf("failed to set PATH environment variable: %w", err)
	}

	// Add to activity state
	return m.state.AddPath(normalizedPath, global)
}

// RemovePath removes a path from the PATH environment variable
func (m *EnvService) RemovePath(path string, global bool) error {
	// Normalize the path
	normalizedPath := util.NormalizePath(path)

	// Remove from session PATH
	currentPath := os.Getenv("PATH")
	pathList := util.SplitPathList(currentPath)

	found := false
	newPathList := []string{}

	for _, p := range pathList {
		if p != normalizedPath {
			newPathList = append(newPathList, p)
		} else {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("path not found in PATH: %s", normalizedPath)
	}

	// Update PATH environment variable
	newPath := util.JoinPathList(newPathList)
	if err := os.Setenv("PATH", newPath); err != nil {
		return fmt.Errorf("failed to set PATH environment variable: %w", err)
	}

	// Remove from activity state
	return m.state.RemovePath(normalizedPath, global)
}

// ListPaths lists PATH entries
func (m *EnvService) ListPaths() []models.PathEntry {
	var paths []models.PathEntry
	for _, entry := range m.state.Global().ActivePaths {
		paths = append(paths, models.PathEntry{
			Path:     entry,
			Priority: 0,
			IsActive: true,
			Scope:    "global",
		})
	}

	for _, entry := range m.state.Session().ActivePaths {
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
func (m *EnvService) SearchPath(path string) []string {
	normalizedPath := util.NormalizePath(path)
	var matches []string

	// Search in active paths
	for _, p := range m.state.Global().ActivePaths {
		if strings.Contains(p, normalizedPath) {
			matches = append(matches, p)
		}
	}

	// Also search in current system PATH
	currentPath := os.Getenv("PATH")
	pathList := util.SplitPathList(currentPath)
	for _, p := range pathList {
		if strings.Contains(p, normalizedPath) {
			matches = append(matches, p)
		}
	}

	return matches
}
