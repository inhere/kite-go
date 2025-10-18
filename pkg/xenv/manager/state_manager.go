package manager

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gookit/goutil"
	"github.com/inhere/kite-go/pkg/xenv/models"
)

// StateManager manages the state data of the environment
type StateManager struct {
	init      bool
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

//
// Tool state management
//

// ActivateTool activates a specific tool version
func (m *StateManager) ActivateTool(name, version string, global bool) error {
	m.ensureInit()

	if global {
		m.global.ActiveTools[name] = version
		return m.SaveGlobalState()
	}

	m.session.ActiveTools[name] = version
	return nil
}

// DeactivateTool deactivates a specific tool version
func (m *StateManager) DeactivateTool(name, version string, global bool) error {
	m.ensureInit()
	if global {
		// Check if the tool is currently active
		currentVersion, exists := m.global.ActiveTools[name]
		if !exists || currentVersion != version {
			return fmt.Errorf("tool %s:%s is not currently active", name, version)
		}

		delete(m.global.ActiveTools, name)
		return m.SaveGlobalState()
	}

	// Check if the tool is currently active
	currentVersion, exists := m.session.ActiveTools[name]
	if !exists || currentVersion != version {
		return fmt.Errorf("tool %s:%s is not currently active", name, version)
	}
	delete(m.session.ActiveTools, name)
	return nil
}

//
// Env state management
//

// SetEnv sets an environment variable
func (m *StateManager) SetEnv(name, value string, global bool) error {
	m.ensureInit()

	// Set global env
	if global {
		m.global.ActiveEnv[name] = value
		return m.SaveGlobalState()
	}

	m.session.ActiveEnv[name] = value
	return nil
}

// UnsetEnv unsets an environment variable
func (m *StateManager) UnsetEnv(name string, global bool) error {
	m.ensureInit()

	// Unset global env
	if global {
		delete(m.global.ActiveEnv, name)
		return m.SaveGlobalState()
	}

	delete(m.session.ActiveEnv, name)
	return nil
}

// AddPath adds a path to the PATH environment variable
func (m *StateManager) AddPath(path string, global bool) error {
	m.ensureInit()

	if global {
		m.global.ActivePaths = append(m.global.ActivePaths, path)
		return m.SaveGlobalState()
	}

	m.session.ActivePaths = append(m.session.ActivePaths, path)
	return nil
}

// RemovePath removes a path from the PATH environment variable
func (m *StateManager) RemovePath(path string, global bool) error {
	m.ensureInit()

	if global {
		newPaths := make([]string, 0, len(m.global.ActivePaths))
		for _, p := range m.global.ActivePaths {
			if p != path {
				newPaths = append(newPaths, p)
			}
		}
		m.global.ActivePaths = newPaths
		return m.SaveGlobalState()
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

//
// Helper methods
//

// Global returns the global activity state
func (m *StateManager) Global() *models.ActivityState {
	return m.global
}

// Session returns the session activity state
func (m *StateManager) Session() *models.ActivityState {
	return m.session
}

// LoadGlobalState loads the global state from file
func (m *StateManager) LoadGlobalState() error {
	// Load global activity state
	globalState, err := models.LoadGlobalState()
	if err != nil {
		return err
	}
	m.global = globalState
	return nil
}

// SaveGlobalState saves the global state to file
func (m *StateManager) SaveGlobalState() error {
	return m.global.Save(m.stateFile)
}

func (m *StateManager) ensureInit() {
	if !m.init {
		goutil.PanicErr(m.Init())
	}
}
