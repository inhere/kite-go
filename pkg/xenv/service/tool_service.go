package service

import (
	"fmt"

	"github.com/gookit/goutil/errorx"
	"github.com/gookit/goutil/strutil"
	"github.com/gookit/goutil/x/ccolor"
	"github.com/inhere/kite-go/pkg/xenv/manager"
	"github.com/inhere/kite-go/pkg/xenv/models"
	"github.com/inhere/kite-go/pkg/xenv/tools"
)

// ToolService handles tool chain management operations
type ToolService struct {
	config *models.Configuration
	state   *manager.StateManager
	toolMgr *manager.ToolManager
	// envMgr *manager.EnvManager
}

// NewToolService creates a new ToolService
func NewToolService(config *models.Configuration, state *manager.StateManager, toolMgr *manager.ToolManager) *ToolService {
	return &ToolService{
		config: config,
		state:   state,
		toolMgr: toolMgr,
		// envMgr: manager.NewEnvManager(),
	}
}

func (ts *ToolService) Register(name string, version string, url string, bin string) error {
	return errorx.Raw("TODO register ...")
}

// ListAll lists all tools
func (ts *ToolService) ListAll(showAll bool) error {
	cfgSdks := ts.config.SDKs
	if len(cfgSdks) == 0 {
		fmt.Println("No SDK tools for managed. see config: sdks, tools")
		return nil
	}
	if err := ts.toolMgr.InitLoad(); err != nil {
		return err
	}

	ccolor.Cyanf("Managed SDK Tools(%d):\n", len(cfgSdks))
	// dump.P(ts.toolMgr.LocalIndexes().SDKs)

	for _, toolCfg := range cfgSdks {
		ccolor.Magentaf(" %s", toolCfg.Name)
		if len(toolCfg.Alias) > 0 {
			fmt.Printf("(Alias: %v) SDK:\n", toolCfg.Alias)
		} else {
			fmt.Println(" SDK:")
		}
		fmt.Printf("  - InstallDir: %s\n", toolCfg.InstallDir)
		if len(toolCfg.BinPaths) > 0 {
			fmt.Printf("  - BinPaths: %v\n", toolCfg.BinPaths)
		}

		locals := ts.toolMgr.ListSDKVersions(toolCfg.Name)
		fmt.Print("  - Installed: ")
		if len(locals) > 0 {
			for _, local := range locals {
				ccolor.Infof("%s ", local.Version)
			}
			fmt.Println()
		} else {
			ccolor.Cyanln("None")
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
	for i, tool := range ts.config.SDKs {
		if tool.Name == name {
			if latest == nil {
				// Simple version comparison - in real implementation, we'd use semver
				latest = &ts.config.SDKs[i]
			}
		}
	}
	return latest
}

// endregion
// region Tool Un/Install
//

// InstallTool installs a tool with the specified version
func (ts *ToolService) InstallTool(name, version string) error {
	toolConfig := ts.config.FindSDKConfig(name)
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
	toolConfig := ts.config.FindSDKConfig(name)
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

// endregion
// region SDK Activate
//

// ActivateSDKs activates multiple SDK tools
func (ts *ToolService) ActivateSDKs(useTools []string, opFlag models.OpFlag) (script string, err error) {
	ts.state.SetBatchMode(true)
	defer ts.state.SetBatchMode(false)

	// Generate shell eval scripts
	gen, err1 := getShellGenerator(ts.config)
	if err1 != nil {
		return "", err1
	}

	actParams := models.NewActivateSDKsParams()
	actParams.OpFlag = opFlag

	for _, arg := range useTools {
		// Parse name:version
		spec, err2 := tools.ParseVersionSpec(arg)
		if err2 != nil {
			return "", err2
		}

		// check activate the tool
		localSdk, err3 := ts.checkActivateSDK(spec)
		if err3 != nil {
			return "", fmt.Errorf("failed to activate tool %q: %w", spec, err3)
		}

		// 如果sdk已经激活过，先要删除之前的激活版本设置的 ENV, PATH
		oldActiveVer := ts.state.Merged().SDKs[spec.Name]
		if oldActiveVer != "" {
			if oldActiveVer == localSdk.Version {
				ccolor.Warnf("The tool %s is already activated, please deactivate it first.\n", localSdk.ID)
				continue
			}

			oldSdk := ts.toolMgr.FindSdkByID(spec.Name + ":" + oldActiveVer)
			if oldSdk != nil {
				actParams.AddRemPath(oldSdk.BinDirPath())
			}
		}

		// 存储激活的真实版本
		actParams.AddSdk(spec.Name, localSdk.Version)
		// Add to PATH and ENVs
		actParams.AddPath(localSdk.BinDirPath())
		if len(localSdk.Config.ActiveEnv) > 0 {
			actParams.AddSetEnvs(localSdk.RenderActiveEnv())
		}

		if opFlag == models.OpFlagGlobal {
			ccolor.Infof("Activate %s as global default\n", localSdk.ID)
		} else if opFlag == models.OpFlagDirenv {
			ccolor.Infof("Activate %s for direnv state\n", localSdk.ID)
		} else {
			ccolor.Infof("Activate %s for current session\n", localSdk.ID)
		}
	}

	// 在shell hook环境中, 生成ENV set脚本
	var sb strutil.Builder
	if gen != nil {
		script1 := gen.GenRemThenAddPaths(actParams.RemPaths, actParams.AddPaths)
		sb.Writeln(script1)
		if len(actParams.AddEnvs) > 0 {
			sb.Writeln(gen.GenSetEnvs(actParams.AddEnvs))
		}
	} else {
		ccolor.Warnln("TIP: The operation will not take effect, please setup the SHELL HOOK first.")
	}

	// Update the activity state
	err = ts.state.UseSDKsWithParams(actParams)
	if err != nil {
		return "", err
	}

	err = ts.state.SaveStateFile()
	return sb.String(), err
}

// check for activates a specific tool version
func (ts *ToolService) checkActivateSDK(spec *tools.VersionSpec) (*models.InstalledTool, error) {
	// Check if the tool is definition
	toolCfg := ts.config.FindSDKConfig(spec.Name)
	if toolCfg == nil {
		return nil, fmt.Errorf("tool %s config is not definition", spec.Name)
	}

	localSdks := ts.toolMgr.ListSDKVersions(toolCfg.Name)
	if len(localSdks) == 0 {
		return nil, fmt.Errorf("sdk tool %s is not installed locally", spec.Name)
	}

	// 根据版本匹配本地可用的sdk
	localSdk := ts.toolMgr.MatchSdkByVersion(localSdks, spec.Version)
	if localSdk == nil {
		return nil, fmt.Errorf("sdk tool %s is not installed locally", spec.ID())
	}

	// 绑定配置信息
	localSdk.Config = toolCfg
	return localSdk, nil
}

// endregion
// region SDK Deactivate
//

// DeactivateSDKs deactivates multiple SDK tools at once
func (ts *ToolService) DeactivateSDKs(deTools []string, opFlag models.OpFlag) (script string, err error) {
	ts.state.SetBatchMode(true)
	defer ts.state.SetBatchMode(false)

	// Generate shell eval scripts
	gen, err1 := getShellGenerator(ts.config)
	if err1 != nil {
		return "", err1
	}

	var delPaths, delEnvs []string

	for _, arg := range deTools {
		spec, err2 := tools.ParseVersionSpec(arg)
		if err2 != nil {
			return "", err2
		}

		// Deactivate the tool
		localSdk, err3 := ts.checkDeactivateSDK(spec, opFlag)
		if err3 != nil {
			ccolor.Warnf("WARN: failed to deactivate tool %q: %w", spec, err3)
			continue
		}

		if localSdk != nil {
			delPaths = append(delPaths, localSdk.BinDirPath())
			if len(localSdk.Config.ActiveEnv) > 0 {
				delEnvs = append(delEnvs, localSdk.Config.ActiveEnvNames()...)
			}
		}

		if opFlag == models.OpFlagGlobal {
			ccolor.Infof("Deactivate %s for global stae\n", spec)
		} else if opFlag == models.OpFlagDirenv {
			ccolor.Infof("Deactivate %s for direnv state\n", spec)
		} else {
			ccolor.Infof("Deactivate %s for current session\n", spec)
		}
	}

	// 在shell hook环境中, 生成ENV remove脚本
	var sb strutil.Builder
	if gen != nil && len(delPaths) > 0 {
		script1, notFounds := gen.GenRemovePaths(delPaths)
		if len(notFounds) > 0 {
			ccolor.Warnf("WARN: %d paths not found in PATH: %v\n", len(notFounds), notFounds)
		}

		sb.Writeln(script1)
		if len(delEnvs) > 0 {
			sb.Writeln(gen.GenUnsetEnvs(delEnvs))
		}
	}

	err = ts.state.DelSDKsWithEnvsPaths(deTools, delEnvs, delPaths, opFlag)
	if err != nil {
		return "", err
	}

	err = ts.state.SaveStateFile()
	return sb.String(), err
}

// deactivateTool deactivates a specific tool version
func (ts *ToolService) checkDeactivateSDK(spec *tools.VersionSpec, opFlag models.OpFlag) (*models.InstalledTool, error) {
	// Check if the tool is definition
	toolCfg := ts.config.FindSDKConfig(spec.Name)
	if toolCfg == nil {
		return nil, fmt.Errorf("sdk %s config is not definition", spec.Name)
	}

	localSdks := ts.toolMgr.ListSDKVersions(toolCfg.Name)
	if len(localSdks) == 0 {
		return nil, fmt.Errorf("sdk %s is not installed locally", spec.Name)
	}

	// 根据版本匹配本地可用的sdk
	localSdk := ts.toolMgr.MatchSdkByVersion(localSdks, spec.Version)
	if localSdk == nil {
		// installed := strings.Join(arrutil.Map1(localSdks, func(t models.InstalledTool) string {
		// 	return t.ID
		// }), ", ")
		// return nil, fmt.Errorf("sdk %s is not installed locally(installed: %s)", spec.ID(), installed)
		return nil, nil
	}

	// Update the activity state
	localSdk.Config = toolCfg
	return localSdk, nil
}
