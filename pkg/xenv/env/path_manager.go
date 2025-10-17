package env

import (
	"fmt"
	"os"
	"strings"

	"github.com/inhere/kite-go/internal/util"
	"github.com/inhere/kite-go/pkg/xenv/models"
)

// PathManager handles PATH environment variable management
type PathManager struct {
	config        *models.Configuration
	activityState *models.ActivityState
}

// NewPathManager creates a new PathManager
func NewPathManager(config *models.Configuration, activityState *models.ActivityState) *PathManager {
	return &PathManager{
		config:        config,
		activityState: activityState,
	}
}

// AddPath adds a path to the PATH environment variable
func (pm *PathManager) AddPath(path string, global bool) error {
	// Normalize the path
	normalizedPath := util.NormalizePath(path)

	// Check if path exists
	if _, err := os.Stat(normalizedPath); os.IsNotExist(err) {
		return fmt.Errorf("path does not exist: %s", normalizedPath)
	}

	if global {
		// Add to global configuration
		// Check if path already exists in global paths
		for _, p := range pm.config.GlobalPaths {
			if p.Path == normalizedPath {
				return fmt.Errorf("path already exists in global paths: %s", normalizedPath)
			}
		}

		// Add the path to global paths
		newPathEntry := models.PathEntry{
			Path:     normalizedPath,
			Scope:    "global",
			IsActive: true,
		}
		pm.config.GlobalPaths = append(pm.config.GlobalPaths, newPathEntry)
	} else {
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
		pm.activityState.ActivePaths = append(pm.activityState.ActivePaths, normalizedPath)
	}

	return nil
}

// RemovePath removes a path from the PATH environment variable
func (pm *PathManager) RemovePath(path string, global bool) error {
	// Normalize the path
	normalizedPath := util.NormalizePath(path)

	if global {
		// Remove from global configuration
		found := false
		newGlobalPaths := []models.PathEntry{}
		
		for _, p := range pm.config.GlobalPaths {
			if p.Path != normalizedPath {
				newGlobalPaths = append(newGlobalPaths, p)
			} else {
				found = true
			}
		}
		
		if !found {
			return fmt.Errorf("path not found in global paths: %s", normalizedPath)
		}
		
		pm.config.GlobalPaths = newGlobalPaths
	} else {
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
		for _, p := range pm.activityState.ActivePaths {
			if p != normalizedPath {
				newActivePaths = append(newActivePaths, p)
			}
		}
		pm.activityState.ActivePaths = newActivePaths
	}

	return nil
}

// ListPaths lists PATH entries
func (pm *PathManager) ListPaths() []models.PathEntry {
	// Return the global PATH entries
	return pm.config.GlobalPaths
}

// SearchPath searches for a path in PATH
func (pm *PathManager) SearchPath(path string) []string {
	normalizedPath := util.NormalizePath(path)
	var matches []string
	
	// Search in global paths
	for _, p := range pm.config.GlobalPaths {
		if strings.Contains(p.Path, normalizedPath) {
			matches = append(matches, p.Path)
		}
	}
	
	// Search in active paths (session)
	for _, p := range pm.activityState.ActivePaths {
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