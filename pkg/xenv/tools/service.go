package tools

import (
	"fmt"
	"path/filepath"

	"github.com/inhere/kite-go/internal/util"
	"github.com/inhere/kite-go/pkg/xenv/models"
)

// ToolService handles tool chain management operations
type ToolService struct {
	config *models.Configuration
}

// NewToolService creates a new ToolService
func NewToolService(config *models.Configuration) *ToolService {
	return &ToolService{
		config: config,
	}
}

// InstallTool installs a tool with the specified version
func (ts *ToolService) InstallTool(name, version string) error {
	// Check if tool is already installed
	id := fmt.Sprintf("%s:%s", name, version)
	for _, tool := range ts.config.Tools {
		if tool.ID == id {
			return fmt.Errorf("tool %s is already installed", id)
		}
	}

	// Prepare installation directory
	installPath := filepath.Join(util.ExpandHome(ts.config.InstallDir), name, version)
	if err := util.EnsureDir(installPath); err != nil {
		return fmt.Errorf("failed to create installation directory: %w", err)
	}

	// Create a new ToolChain instance
	toolChain := models.ToolChain{
		ID:        id,
		Name:      name,
		Version:   version,
		InstallDir: installPath,
		Installed: true,
		BinPaths:  []string{installPath}, // Default to install directory
	}

	// In a real implementation, we would download and install the tool here
	downloader := NewDownloader()
	err := downloader.DownloadFile(toolChain.InstallURL, ts.config.DownloadDir)
	if err != nil {
		return err
	}

	// TODO save tool to local.json

	return nil
}

// UninstallTool uninstalls a tool with the specified version
func (ts *ToolService) UninstallTool(name, version string) error {
	id := fmt.Sprintf("%s:%s", name, version)

	// Find the tool in the configuration
	foundIndex := -1
	for i, tool := range ts.config.Tools {
		if tool.ID == id {
			foundIndex = i
			break
		}
	}

	if foundIndex == -1 {
		return fmt.Errorf("tool %s is not installed", id)
	}

	// TODO Remove the tool from local.json

	return nil
}

// ListTools returns all managed tools
func (ts *ToolService) ListTools() []models.ToolChain {
	return ts.config.Tools
}

// UpdateTool updates a tool to the specified version
func (ts *ToolService) UpdateTool(name, version string) error {
	// For update, we'll install the new version
	return ts.InstallTool(name, version)
}

// GetTool returns information about a specific tool
func (ts *ToolService) GetTool(name string) *models.ToolChain {
	// Find the latest version of the tool
	var latest *models.ToolChain
	for i, tool := range ts.config.Tools {
		if tool.Name == name {
			if latest == nil || tool.Version > latest.Version {
				// Simple version comparison - in real implementation, we'd use semver
				latest = &ts.config.Tools[i]
			}
		}
	}
	return latest
}

// EnsureBinDir ensures the bin directory exists and is in the PATH
func (ts *ToolService) EnsureBinDir() error {
	binDir := util.ExpandHome(ts.config.BinDir)
	return util.EnsureDir(binDir)
}
