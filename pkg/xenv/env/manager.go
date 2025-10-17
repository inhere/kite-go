package env

import (
	"fmt"
	"os"

	"github.com/inhere/kite-go/pkg/xenv/models"
)

// Manager handles environment variable and PATH management
type Manager struct {
	config        *models.Configuration
	activityState *models.ActivityState
}

// NewManager creates a new Manager
func NewManager(config *models.Configuration, activityState *models.ActivityState) *Manager {
	return &Manager{
		config:        config,
		activityState: activityState,
	}
}

// SetEnv sets an environment variable
func (m *Manager) SetEnv(name, value string, global bool) error {
	if global {
		// Add to global configuration
		if m.config.GlobalEnv == nil {
			m.config.GlobalEnv = make(map[string]models.EnvVariable)
		}
		m.config.GlobalEnv[name] = models.EnvVariable{
			Name:     name,
			Value:    value,
			Scope:    "global",
			IsActive: true,
		}
	} else {
		// Add to session (would be handled differently in practice)
		// For now, we'll just validate the variable can be set
		if err := os.Setenv(name, value); err != nil {
			return fmt.Errorf("failed to set environment variable: %w", err)
		}

		// Add to activity state to track session variables
		if m.activityState.ActiveEnv == nil {
			m.activityState.ActiveEnv = make(map[string]string)
		}
		m.activityState.ActiveEnv[name] = value
	}

	return nil
}

// UnsetEnv unsets an environment variable
func (m *Manager) UnsetEnv(name string, global bool) error {
	if global {
		// Remove from global configuration
		delete(m.config.GlobalEnv, name)
	} else {
		// Remove from session
		if err := os.Unsetenv(name); err != nil {
			return fmt.Errorf("failed to unset environment variable: %w", err)
		}

		// Remove from activity state
		delete(m.activityState.ActiveEnv, name)
	}

	return nil
}

// ListEnv lists environment variables
func (m *Manager) ListEnv() map[string]models.EnvVariable {
	// Return the global environment variables
	return m.config.GlobalEnv
}
