package tools

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/inhere/kite-go/internal/util"
	"github.com/inhere/kite-go/pkg/xenv/models"
)

// Activator handles tool chain activation
type Activator struct {
	config       *models.Configuration
	activityState *models.ActivityState
}

// NewActivator creates a new Activator
func NewActivator(config *models.Configuration, activityState *models.ActivityState) *Activator {
	return &Activator{
		config:        config,
		activityState: activityState,
	}
}

// ActivateTool activates a specific tool version
func (a *Activator) ActivateTool(name, version string, global bool) error {
	id := fmt.Sprintf("%s:%s", name, version)

	// Check if the tool is installed
	toolFound := false
	for _, tool := range a.config.Tools {
		if tool.ID == id {
			toolFound = true
			break
		}
	}

	if !toolFound {
		return fmt.Errorf("tool %s is not installed", id)
	}

	// Update the activity state
	if a.activityState.ActiveTools == nil {
		a.activityState.ActiveTools = make(map[string]string)
	}
	
	a.activityState.ActiveTools[name] = version

	// If global flag is set, we might want to persist this in configuration
	if global {
		// For global activation, save the state
		// In a real implementation, this would involve writing to a global state file
		if err := a.saveGlobalState(); err != nil {
			return fmt.Errorf("failed to save global state: %w", err)
		}
	}

	return nil
}

// saveGlobalState saves the global activation state to a file
func (a *Activator) saveGlobalState() error {
	// Get the config directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, ".config", "xenv")
	if err := util.EnsureDir(configDir); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// In a real implementation, we would write the actual activity state to the file
	// For now, this is a placeholder

	return nil
}