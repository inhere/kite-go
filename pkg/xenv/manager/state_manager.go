package manager

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gookit/goutil"
	"github.com/gookit/goutil/fsutil"
	"github.com/inhere/kite-go/pkg/xenv/models"
)

// StateManager manages the state data of the environment
type StateManager struct {
	init      bool
	batchMode bool
	stateFile string
	// global state data
	global  *models.ActivityState
	session *models.ActivityState
}

func stateFilePath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	return filepath.Join(homeDir, ".config", "xenv", "activity.json")
}

// NewStateManager creates a new StateManager
func NewStateManager() *StateManager {
	return &StateManager{
		stateFile: stateFilePath(),
		session:   models.NewActivityState(),
	}
}

// Init initializes the state manager
func (m *StateManager) Init() error {
	if m.init {
		return nil
	}
	m.init = true
	return m.LoadGlobalState()
}

// StateFile returns the state file path
func (m *StateManager) StateFile() string {
	return m.stateFile
}

// SetBatchMode sets the batch mode flag
func (m *StateManager) SetBatchMode(enabled bool) {
	m.batchMode = enabled
}

// endregion
// region Tool state management
//

// UseToolsWithEnvsPaths activates multiple tools and with envs, paths
func (m *StateManager) UseToolsWithEnvsPaths(tools, envs map[string]string, paths []string, global bool) error {
	if err := m.Init(); err != nil {
		return err
	}

	if global {
		m.global.AddToolsWithEnvsPaths(tools, envs, paths)
		// Save the global state
		if !m.batchMode {
			return m.SaveGlobalState()
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
		m.global.ActiveTools[name] = version
		// Save the global state
		if !m.batchMode {
			return m.SaveGlobalState()
		}
	} else {
		m.session.ActiveTools[name] = version
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
			return m.SaveGlobalState()
		}
	} else {
		m.session.DelToolsWithEnvsPaths(tools, envs, paths)
	}
	return nil
}

// DeactivateTool deactivates a specific tool version
func (m *StateManager) DeactivateTool(name, version string, global bool) error {
	if err := m.Init(); err != nil {
		return err
	}

	if global {
		// Check if the tool is currently active
		currentVersion, exists := m.global.ActiveTools[name]
		if !exists || currentVersion != version {
			return fmt.Errorf("tool %s:%s is not currently active", name, version)
		}

		delete(m.global.ActiveTools, name)
		// Save the global state
		if !m.batchMode {
			return m.SaveGlobalState()
		}
		return nil
	}

	// Check if the tool is currently active
	currentVersion, exists := m.session.ActiveTools[name]
	if !exists || currentVersion != version {
		return fmt.Errorf("tool %s:%s is not currently active", name, version)
	}
	delete(m.session.ActiveTools, name)
	return nil
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
			m.global.ActiveEnv[name] = value
		}
		// Save the global state
		if !m.batchMode {
			return m.SaveGlobalState()
		}
		return nil
	}

	for name, value := range envs {
		m.session.ActiveEnv[name] = value
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
		m.global.ActiveEnv[name] = value
		// Save the global state
		if !m.batchMode {
			return m.SaveGlobalState()
		}
	} else {
		// Set session env
		m.session.ActiveEnv[name] = value
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
		if _, exists := m.global.ActiveEnv[name]; exists {
			// return fmt.Errorf("environment variable %s is not currently set", name)
			delete(m.global.ActiveEnv, name)
			// Save the global state
			if !m.batchMode {
				return m.SaveGlobalState()
			}
		}
	} else {
		// check exists
		if _, exists := m.session.ActiveEnv[name]; exists {
			// return fmt.Errorf("environment variable %s is not currently set", name)
			delete(m.session.ActiveEnv, name)
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
			return m.SaveGlobalState()
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
			return m.SaveGlobalState()
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
		newPaths := make([]string, 0, len(m.global.ActivePaths))
		for _, p := range m.global.ActivePaths {
			if p != path {
				newPaths = append(newPaths, p)
			}
		}
		m.global.ActivePaths = newPaths
		// Save the global state
		if !m.batchMode {
			return m.SaveGlobalState()
		}
	}

	newPaths := make([]string, 0, len(m.session.ActivePaths))
	for _, p := range m.session.ActivePaths {
		if p != path {
			newPaths = append(newPaths, p)
		}
	}
	m.session.ActivePaths = newPaths
	return nil
}

// endregion
// region Load/Save GlobalState
//

// LoadGlobalState loads the global state from file
func (m *StateManager) LoadGlobalState() error {
	// Check if file exists
	if _, err := os.Stat(m.stateFile); os.IsNotExist(err) {
		// Return default state if file doesn't exist
		m.global = models.NewActivityState()
		return nil
	}

	// Read the file
	data, err := os.ReadFile(m.stateFile)
	if err != nil {
		return err
	}

	// Unmarshal the JSON
	var state models.ActivityState
	if err1 := json.Unmarshal(data, &state); err1 != nil {
		return err1
	}

	m.global = &state
	return nil
}

// SaveGlobalState saves the global state to file
func (m *StateManager) SaveGlobalState() error {
	if err := fsutil.MkParentDir(m.stateFile); err != nil {
		return err
	}

	// Update timestamps
	m.global.UpdatedAt = time.Now()
	if m.global.CreatedAt.IsZero() {
		m.global.CreatedAt = m.global.UpdatedAt
	}

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

// Global returns the global activity state
func (m *StateManager) Global() *models.ActivityState {
	return m.global
}

// Session returns the session activity state
func (m *StateManager) Session() *models.ActivityState {
	return m.session
}

func (m *StateManager) requireInit() {
	goutil.PanicErr(m.Init())
}
