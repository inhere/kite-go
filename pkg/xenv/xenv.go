package xenv

import (
	"fmt"

	"github.com/inhere/kite-go/pkg/xenv/config"
	"github.com/inhere/kite-go/pkg/xenv/manager"
	"github.com/inhere/kite-go/pkg/xenv/service"
)

// Init initializes the xenv config and state
func Init() error {
	// Initialize configuration
	if err := config.Mgr.Init(); err != nil {
		return fmt.Errorf("failed to initialize configuration: %w", err)
	}

	if err := stateMgr.Init(); err != nil {
		return fmt.Errorf("failed to load global state: %w", err)
	}
	return nil
}

var stateMgr = manager.NewStateManager()

// State returns the state manager
func State() *manager.StateManager {
	return stateMgr
}

func EnvService() (*service.EnvService, error) {
	// Initialize configuration
	if err := config.Mgr.Init(); err != nil {
		return nil, fmt.Errorf("failed to initialize configuration: %w", err)
	}

	// Create env manager
	return service.NewEnvService(config.Mgr.Config, stateMgr), nil
}

