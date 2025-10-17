package tools

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/inhere/kite-go/internal/util"
	"github.com/inhere/kite-go/pkg/xenv/models"
)

// Uninstaller handles uninstalling tools
type Uninstaller struct {
	service *ToolService
	config  *models.Configuration
}

// NewUninstaller creates a new Uninstaller
func NewUninstaller(service *ToolService, config *models.Configuration) *Uninstaller {
	return &Uninstaller{
		service: service,
		config:  config,
	}
}

// Uninstall removes a tool with the specified name and version
func (u *Uninstaller) Uninstall(name, version string, keepConfig bool) error {
	id := fmt.Sprintf("%s:%s", name, version)
	
	// Find the tool in the configuration
	var tool *models.ToolChain
	var toolIndex int
	found := false
	
	for i, t := range u.config.Tools {
		if t.ID == id {
			tool = &u.config.Tools[i]
			toolIndex = i
			found = true
			break
		}
	}
	
	if !found {
		return fmt.Errorf("tool %s is not installed", id)
	}

	// Remove the tool from the bin directory (remove shims)
	if err := u.removeShims(tool); err != nil {
		// Continue execution even if removing shims fails
	}

	// Remove the installed tool directory if not keeping config
	if !keepConfig {
		if err := os.RemoveAll(tool.InstallDir); err != nil {
			return fmt.Errorf("failed to remove installation directory: %w", err)
		}
	} 

	// Remove the tool from configuration
	u.config.Tools = append(u.config.Tools[:toolIndex], u.config.Tools[toolIndex+1:]...)

	return nil
}

// removeShims removes the symlinks (shims) for the tool executables
func (u *Uninstaller) removeShims(tool *models.ToolChain) error {
	binDir := util.ExpandHome(u.service.config.BinDir)
	
	// For each binary path of the tool, remove the shim
	for _, binPath := range tool.BinPaths {
		// Get all executable files in the bin path
		entries, err := os.ReadDir(binPath)
		if err != nil {
			continue // Skip if directory doesn't exist
		}
		
		for _, entry := range entries {
			if !entry.IsDir() && isUninstallExecutable(entry.Name()) {
				// Construct the expected shim path
				shimPath := filepath.Join(binDir, entry.Name())
				
				// Check if the shim exists and remove it
				if _, err := os.Stat(shimPath); err == nil {
					if err := os.Remove(shimPath); err != nil {
						return fmt.Errorf("failed to remove shim %s: %w", shimPath, err)
					}
				}
			}
		}
	}
	
	return nil
}

// isUninstallExecutable checks if a file name suggests it's executable
func isUninstallExecutable(filename string) bool {
	// This is a simplified check - in reality, we'd need to check the actual file
	// For now, we'll assume executables don't have common non-executable extensions
	nonExecutableExts := []string{".txt", ".log", ".md", ".json", ".yaml", ".yml", ".toml", ".xml", ".html", ".css", ".js"}

	for _, ext := range nonExecutableExts {
		if len(filename) > len(ext) && filename[len(filename)-len(ext):] == ext {
			return false
		}
	}

	return true
}