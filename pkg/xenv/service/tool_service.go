package service

import (
	"fmt"

	"github.com/gookit/goutil/errorx"
	"github.com/gookit/goutil/x/ccolor"
	"github.com/inhere/kite-go/internal/util"
	"github.com/inhere/kite-go/pkg/xenv/manager"
	"github.com/inhere/kite-go/pkg/xenv/models"
	"github.com/inhere/kite-go/pkg/xenv/tools"
)

// ToolService handles tool chain management operations
type ToolService struct {
	config *models.Configuration
	state   *manager.StateManager
	toolMgr *manager.ToolManager
}

// NewToolService creates a new ToolService
func NewToolService(config *models.Configuration, state *manager.StateManager, toolMgr *manager.ToolManager) *ToolService {
	return &ToolService{
		config: config,
		state:   state,
		toolMgr: toolMgr,
	}
}

func (ts *ToolService) Register(name string, version string, url string, bin string) error {
	return errorx.Raw("TODO register ...")
}

// InstallTool installs a tool with the specified version
func (ts *ToolService) InstallTool(name, version string) error {
	toolConfig := ts.config.FindToolConfig(name)
	// Check if tool is defined
	if toolConfig == nil {
		return fmt.Errorf("tool %s is not defined in config", name)
	}

	// 查找 local.json 是否存在
	id := fmt.Sprintf("%s:%s", name, version)
	if ts.toolMgr.FindSdkByID(id) != nil {
		return fmt.Errorf("tool %s is already installed in local", id)
	}

	// download and install the tool here
	installer := tools.NewInstaller(ts.config)
	err := installer.Install(toolConfig, version)
	if err != nil {
		return err
	}

	// save tool to local.json
	return ts.toolMgr.AddSDKTool(name, version, installer.InstallDir)
}

// Uninstall uninstalls a sdk tool with the specified version
func (ts *ToolService) Uninstall(name, version string) error {
	id := fmt.Sprintf("%s:%s", name, version)

	// Find the tool in the configuration
	toolConfig := ts.config.FindToolConfig(name)
	if toolConfig == nil {
		return fmt.Errorf("tool %s is not installed", id)
	}

	// TODO 从 state 里检测并删除

	// 查找 local.json 是否存在
	localTool := ts.toolMgr.FindSdkByID(id)
	if localTool == nil {
		return fmt.Errorf("tool %s:%s is not installed", name, version)
	}

	uninstaller := tools.NewUninstaller(ts.config)
	err := uninstaller.Uninstall(toolConfig, localTool, false)
	if err != nil {
		return err
	}

	// remove from ts.localTools and save local.json
	return ts.toolMgr.DeleteSDKTool(localTool)
}

// ListAll lists all tools
func (ts *ToolService) ListAll(showAll bool) error {
	cfgTools := ts.config.Tools
	if len(cfgTools) == 0 {
		fmt.Println("No tools for managed. see config: tools, simple_tools")
		return nil
	}

	ccolor.Cyanf("Managed SDK Tools(%d):\n", len(cfgTools))

	for _, toolCfg := range cfgTools {
		ccolor.Magentaf("%s", toolCfg.Name)
		if len(toolCfg.Alias) > 0 {
			fmt.Printf("(Aliases: %v) SDK:\n", toolCfg.Alias)
		} else {
			fmt.Println(" SDK:")
		}
		fmt.Printf("  - InstallDir: %s\n", toolCfg.InstallDir)
		if len(toolCfg.BinPaths) > 0 {
			fmt.Printf("  - BinPaths: %v\n", toolCfg.BinPaths)
		}

		locals := ts.toolMgr.ListLocalVersions(toolCfg.Name)
		if len(locals) > 0 {
			fmt.Print("  - Installed: ")
			for _, local := range locals {
				ccolor.Infof("%s ", local.Version)
			}
			fmt.Println()
		}
	}
	return nil
}

func (ts *ToolService) IndexLocalTools() error {
	return ts.toolMgr.IndexLocalTools()
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

//
// Tool Activation
//

// ActivateTools activates multiple tools
func (ts *ToolService) ActivateTools(useTools []string, global bool) error {
	ts.state.SetBatchMode(true)
	defer ts.state.SetBatchMode(false)

	for _, arg := range useTools {
		// Parse name:version
		spec, err := tools.ParseVersionSpec(arg)
		if err != nil {
			return err
		}

		// Activate the tool
		if err := ts.activateTool(spec, global); err != nil {
			return fmt.Errorf("failed to activate tool %q: %w", spec, err)
		}

		if global {
			ccolor.Infof("Activate %s as global default\n", spec)
		} else {
			ccolor.Infof("Activate %s for current session\n", spec)
		}
	}

	if global {
		ccolor.Printf("Global: save state to %s\n", ts.state.StateFile())
	}
	return ts.state.SaveGlobalState()
}

// activateTool activates a specific tool version
func (ts *ToolService) activateTool(spec *tools.VersionSpec, global bool) error {
	// Check if the tool is definition
	if !ts.config.IsToolDefined(spec.Name) {
		return fmt.Errorf("tool %s config is not definition", spec.Name)
	}

	// Update the activity state
	return ts.state.ActivateTool(spec.Name, spec.Version, global)
}

// DeactivateTools deactivates multiple tools at once
func (ts *ToolService) DeactivateTools(deTools []string, global bool) error {
	ts.state.SetBatchMode(true)
	defer ts.state.SetBatchMode(false)

	for _, arg := range deTools {
		spec, err := tools.ParseVersionSpec(arg)
		if err != nil {
			return err
		}
		if err := ts.deactivateTool(spec.Name, spec.Version, global); err != nil {
			return fmt.Errorf("failed to deactivate tool %q: %w", spec, err)
		}
		if global {
			ccolor.Printf("Deactivate %s for global default\n", spec)
		} else {
			ccolor.Printf("Deactivate %s for current session\n", spec)
		}
	}

	if global {
		ccolor.Printf("Global: save state to %s\n", ts.state.StateFile())
	}
	return ts.state.SaveGlobalState()
}

// deactivateTool deactivates a specific tool version
func (ts *ToolService) deactivateTool(name, version string, global bool) error {
	// Check if the tool is definition
	if !ts.config.IsToolDefined(name) {
		return fmt.Errorf("tool %s config is not definition", name)
	}

	return ts.state.DeactivateTool(name, version, global)
}
