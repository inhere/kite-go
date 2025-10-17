package tools

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/inhere/kite-go/internal/util"
	"github.com/inhere/kite-go/pkg/xenv/models"
)

// Deactivator handles tool chain deactivation
type Deactivator struct {
	config       *models.Configuration
	activityState *models.ActivityState
}

// NewDeactivator creates a new Deactivator
func NewDeactivator(config *models.Configuration, activityState *models.ActivityState) *Deactivator {
	return &Deactivator{
		config:        config,
		activityState: activityState,
	}
}

// DeactivateTool deactivates a specific tool version
func (d *Deactivator) DeactivateTool(name, version string, global bool) error {
	id := fmt.Sprintf("%s:%s", name, version)

	// Check if the tool is currently active
	currentVersion, exists := d.activityState.ActiveTools[name]
	if !exists || currentVersion != version {
		return fmt.Errorf("tool %s is not currently active", id)
	}

	// Remove from active tools
	delete(d.activityState.ActiveTools, name)

	// If global flag is set, we might want to persist this in configuration
	if global {
		// For global deactivation, save the state
		// In a real implementation, this would involve writing to a global state file
		if err := d.saveGlobalState(); err != nil {
			return fmt.Errorf("failed to save global state: %w", err)
		}
	}

	return nil
}

// saveGlobalState saves the global deactivation state to a file
func (d *Deactivator) saveGlobalState() error {
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