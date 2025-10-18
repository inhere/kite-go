package tools

import (
	"fmt"

	"github.com/inhere/kite-go/pkg/xenv/manager"
	"github.com/inhere/kite-go/pkg/xenv/models"
)

// Activator handles tool chain activation
type Activator struct {
	config       *models.Configuration
	state *manager.StateManager
}

// NewActivator creates a new Activator
func NewActivator(config *models.Configuration, state *manager.StateManager) *Activator {
	return &Activator{
		config:        config,
		state: state,
	}
}

// ActivateTool activates a specific tool version
func (a *Activator) ActivateTool(name, version string, global bool) error {
	// Check if the tool is definition
	if !a.config.IsToolDefined(name) {
		return fmt.Errorf("tool %s:%s config is not definition", name, version)
	}

	// Update the activity state
	return a.state.ActivateTool(name, version, global)
}

// DeactivateTool deactivates a specific tool version
func (a *Activator) DeactivateTool(name, version string, global bool) error {
	// Check if the tool is definition
	if !a.config.IsToolDefined(name) {
		return fmt.Errorf("tool %s:%s config is not definition", name, version)
	}

	return a.state.DeactivateTool(name, version, global)
}

// saveGlobalState saves the global activation state to a file
func (a *Activator) saveGlobalState() error {
	// In a real implementation, we would write the actual activity state to the file
	// For now, this is a placeholder
	return a.state.SaveGlobalState()
}
