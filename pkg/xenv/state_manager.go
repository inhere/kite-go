package xenv

import "github.com/inhere/kite-go/pkg/xenv/models"

// StateManager manages the state data of the environment
type StateManager struct {
	config  *models.Configuration
	global  *models.ActivityState
	session *models.ActivityState
}

// NewStateManager creates a new StateManager
func NewStateManager(config *models.Configuration, globalState *models.ActivityState) *StateManager {
	return &StateManager{
		config:  config,
		global:  globalState,
		session: &models.ActivityState{},
	}
}

func (sm *StateManager) LoadState() error {
	// Load global activity state
	globalState, err := models.LoadGlobalState()
	if err != nil {
		return err
	}
	sm.global = globalState
	return nil
}

// SaveState saves the global state to file
func (sm *StateManager) SaveState() error {
	return sm.global.Save(sm.config.ConfigDir())
}
