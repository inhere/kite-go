package tools

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/jsonutil"
	"github.com/inhere/kite-go/internal/util"
	"github.com/inhere/kite-go/pkg/xenv/models"
)

// ToolService handles tool chain management operations
type ToolService struct {
	config      *models.Configuration
	loaded      bool
	metaFile    string
	localTools  *models.LocalTools // TODO 安装后
	globalState *models.ActivityState
}

// NewToolService creates a new ToolService
func NewToolService(config *models.Configuration) *ToolService {
	return &ToolService{
		config:     config,
		localTools: &models.LocalTools{Version: "v1"},
	}
}

// LoadLocalTools local installed tools information
func (ts *ToolService) LoadLocalTools() error {
	if ts.loaded {
		return nil
	}
	ts.loaded = true

	ts.metaFile = ts.config.InstallDir + "/local.json"
	if fsutil.IsFile(ts.metaFile) {
		err := jsonutil.DecodeFile(ts.metaFile, ts.localTools)
		if err != nil {
			return err
		}
	}
	return nil
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
		InstallDir: installPath,
		Installed: true,
		BinPaths:  []string{installPath}, // Default to install directory
	}

	// download and install the tool here
	installer := NewInstaller(ts.config)
	err := installer.Install(&toolChain, version)
	if err != nil {
		return err
	}

	// build local installed tool info
	currentTime := time.Now()
	if ts.localTools.CreatedAt.IsZero() {
		ts.localTools.CreatedAt = currentTime
	}
	ts.localTools.UpdatedAt = currentTime
	ts.localTools.SdkTools = append(ts.localTools.SdkTools, models.InstalledTool{
		ID:         id,
		Name:       name,
		Version:    version,
		InstallDir: installPath,
		BinDir:     "",
		BinPaths:   []string{installPath},
		CreatedAt:  currentTime,
		UpdatedAt:  currentTime,
	})

	// TODO save tool to local.json
	return ts.SaveLocalTools()
}

// Uninstall uninstalls a sdk tool with the specified version
func (ts *ToolService) Uninstall(name, version string) error {
	id := fmt.Sprintf("%s:%s", name, version)
	sdkTools := ts.localTools.SdkTools

	// Find the tool in the configuration
	foundIndex := -1
	for i, tool := range sdkTools {
		if tool.ID == id {
			foundIndex = i
			break
		}
	}
	if foundIndex == -1 {
		return fmt.Errorf("tool %s is not installed", id)
	}

	uninstaller := NewUninstaller(ts.config)
	err := uninstaller.Uninstall(name, version, false)
	if err != nil {
		return err
	}

	// remove from ts.localTools
	ts.localTools.SdkTools = append(sdkTools[:foundIndex], sdkTools[foundIndex+1:]...)
	// save local.json
	return ts.SaveLocalTools()
}

// ListTools returns all managed tools
func (ts *ToolService) ListTools() []models.ToolChain {
	return ts.config.Tools
}

// LocalTools returns all installed sdk,simple tools
func (ts *ToolService) LocalTools() *models.LocalTools {
	return ts.localTools
}

// InstalledTools returns all installed sdk tools
func (ts *ToolService) InstalledTools() []models.InstalledTool {
	return ts.localTools.SdkTools
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
			if latest == nil {
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

// SaveLocalTools saves the local tools information
func (ts *ToolService) SaveLocalTools() error {
	return jsonutil.WriteFile(ts.metaFile, ts.localTools)
}
