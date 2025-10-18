package tools

import (
	"fmt"

	"github.com/inhere/kite-go/pkg/xenv/models"
)

// Activator handles tool chain activation
type Activator struct {
	config       *models.Configuration
	globalState *models.ActivityState
}

// NewActivator creates a new Activator
func NewActivator(config *models.Configuration, globalState *models.ActivityState) *Activator {
	return &Activator{
		config:        config,
		globalState: globalState,
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
	if a.globalState.ActiveTools == nil {
		a.globalState.ActiveTools = make(map[string]string)
	}

	a.globalState.ActiveTools[name] = version

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

// DeactivateTool deactivates a specific tool version
func (a *Activator) DeactivateTool(name, version string, global bool) error {
	id := fmt.Sprintf("%s:%s", name, version)

	// Check if the tool is currently active
	currentVersion, exists := a.globalState.ActiveTools[name]
	if !exists || currentVersion != version {
		return fmt.Errorf("tool %s is not currently active", id)
	}

	// Remove from active tools
	delete(a.globalState.ActiveTools, name)

	// If global flag is set, we might want to persist this in configuration
	if global {
		// For global deactivation, save the state
		// In a real implementation, this would involve writing to a global state file
		if err := a.saveGlobalState(); err != nil {
			return fmt.Errorf("failed to save global state: %w", err)
		}
	}

	return nil
}

// saveGlobalState saves the global activation state to a file
func (a *Activator) saveGlobalState() error {
	// In a real implementation, we would write the actual activity state to the file
	// For now, this is a placeholder
	return a.globalState.Save(a.config.ConfigDir())
}
