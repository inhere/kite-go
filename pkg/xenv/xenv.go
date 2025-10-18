package xenv

import (
	"fmt"

	"github.com/inhere/kite-go/pkg/xenv/config"
	"github.com/inhere/kite-go/pkg/xenv/env"
)

func EnvManager() (*env.Manager, error) {
	// Initialize configuration
	if err := config.Mgr.Init(); err != nil {
		return nil, fmt.Errorf("failed to initialize configuration: %w", err)
	}

	// Create env manager
	return env.NewManager(config.Mgr.Config, config.Mgr.State),  nil
}
