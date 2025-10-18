package env

import (
	"fmt"
	"os"
	"strings"

	"github.com/inhere/kite-go/internal/util"
	"github.com/inhere/kite-go/pkg/xenv/models"
)

// Manager handles environment variable and PATH management
type Manager struct {
	config       *models.Configuration
	globalState  *models.ActivityState
	sessionState *models.ActivityState
}

// NewManager creates a new Manager
func NewManager(config *models.Configuration, globalState *models.ActivityState) *Manager {
	return &Manager{
		config:       config,
		globalState:  globalState,
		sessionState: &models.ActivityState{},
	}
}

// SetEnv sets an environment variable
func (m *Manager) SetEnv(name, value string) error {
	// Add to session (would be handled differently in practice)
	// For now, we'll just validate the variable can be set
	if err := os.Setenv(name, value); err != nil {
		return fmt.Errorf("failed to set environment variable: %w", err)
	}

	// Add to activity state to track session variables
	if m.globalState.ActiveEnv == nil {
		m.globalState.ActiveEnv = make(map[string]string)
	}
	m.globalState.ActiveEnv[name] = value

	return nil
}

// UnsetEnv unsets an environment variable
func (m *Manager) UnsetEnv(name string) error {
	// Remove from session
	if err := os.Unsetenv(name); err != nil {
		return fmt.Errorf("failed to unset environment variable: %w", err)
	}

	// Remove from activity state
	delete(m.globalState.ActiveEnv, name)
	return nil
}

// ListEnv lists environment variables
func (m *Manager) ListEnv() map[string]string {
	// Return the global environment variables
	return m.config.GlobalEnv
}

//
// PATH management
//

// AddPath adds a path to the PATH environment variable
func (m *Manager) AddPath(path string) error {
	// Normalize the path
	normalizedPath := util.NormalizePath(path)

	// Check if path exists
	if _, err := os.Stat(normalizedPath); os.IsNotExist(err) {
		return fmt.Errorf("path does not exist: %s", normalizedPath)
	}

	// Add to session PATH
	currentPath := os.Getenv("PATH")
	pathList := util.SplitPathList(currentPath)

	// Check if path already exists
	for _, p := range pathList {
		if p == normalizedPath {
			return fmt.Errorf("path already exists in PATH: %s", normalizedPath)
		}
	}

	// Add the path to the beginning of PATH (highest priority)
	newPathList := append([]string{normalizedPath}, pathList...)
	newPath := util.JoinPathList(newPathList)

	// Set the new PATH environment variable
	if err := os.Setenv("PATH", newPath); err != nil {
		return fmt.Errorf("failed to set PATH environment variable: %w", err)
	}

	// Add to activity state to track session paths
	m.globalState.ActivePaths = append(m.globalState.ActivePaths, normalizedPath)
	return nil
}

// RemovePath removes a path from the PATH environment variable
func (m *Manager) RemovePath(path string) error {
	// Normalize the path
	normalizedPath := util.NormalizePath(path)

	// Remove from session PATH
	currentPath := os.Getenv("PATH")
	pathList := util.SplitPathList(currentPath)

	found := false
	newPathList := []string{}

	for _, p := range pathList {
		if p != normalizedPath {
			newPathList = append(newPathList, p)
		} else {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("path not found in PATH: %s", normalizedPath)
	}

	// Update PATH environment variable
	newPath := util.JoinPathList(newPathList)
	if err := os.Setenv("PATH", newPath); err != nil {
		return fmt.Errorf("failed to set PATH environment variable: %w", err)
	}

	// Remove from activity state
	newActivePaths := []string{}
	for _, p := range m.globalState.ActivePaths {
		if p != normalizedPath {
			newActivePaths = append(newActivePaths, p)
		}
	}
	m.globalState.ActivePaths = newActivePaths

	return nil
}

// ListPaths lists PATH entries
func (m *Manager) ListPaths() []models.PathEntry {
	var paths []models.PathEntry
	for _, entry := range m.globalState.ActivePaths {
		paths = append(paths, models.PathEntry{
			Path:     entry,
			Priority: 0,
			IsActive: true,
			Scope:    "global",
		})
	}

	for _, entry := range m.sessionState.ActivePaths {
		paths = append(paths, models.PathEntry{
			Path:     entry,
			Priority: 0,
			IsActive: true,
			Scope:    "session",
		})
	}
	return paths
}

// SearchPath searches for a path in PATH
func (m *Manager) SearchPath(path string) []string {
	normalizedPath := util.NormalizePath(path)
	var matches []string

	// Search in active paths
	for _, p := range m.globalState.ActivePaths {
		if strings.Contains(p, normalizedPath) {
			matches = append(matches, p)
		}
	}

	// Also search in current system PATH
	currentPath := os.Getenv("PATH")
	pathList := util.SplitPathList(currentPath)
	for _, p := range pathList {
		if strings.Contains(p, normalizedPath) {
			matches = append(matches, p)
		}
	}

	return matches
}

// SaveGlobalState saves the global state to file
func (m *Manager) SaveGlobalState() error {
	return m.globalState.Save(m.config.ConfigDir())
}
