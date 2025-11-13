package manager

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/jsonutil"
	"github.com/gookit/goutil/x/ccolor"
	"github.com/inhere/kite-go/pkg/xenv/models"
	"github.com/inhere/kite-go/pkg/xenv/xenvcom"
	"github.com/inhere/kite-go/pkg/xenv/xenvutil"
)

type ToolManager struct {
	init bool
	// config
	config *models.Configuration
	// local data file
	localLoad bool
	localFile  string
	localTools *models.ToolsLocal
	// tools register data - 从 config 配置中初始化 tools/config.json TODO
	configFile string
	configLoad bool
	// caches TODO
	groupSdks map[string][]models.InstalledTool
}

// NewToolManager creates a new ToolManager instance
func NewToolManager() *ToolManager {
	return &ToolManager{
		localTools: &models.ToolsLocal{Version: "v1"},
		groupSdks: make(map[string][]models.InstalledTool),
	}
}

// Init initializes the state manager
func (m *ToolManager) Init(config *models.Configuration) error {
	if m.init {
		return nil
	}
	m.init = true
	m.config = config
	return nil
}

// InitLoad loads the local indexes
func (m *ToolManager) InitLoad() error {
	return m.ensureLocalLoad(false)
}

// InitLoad1 loads the local indexes
// func (m *ToolManager) InitLoad1(cfg *models.Configuration) error {
// 	_ = m.Init(cfg)
// 	return m.ensureLocalLoad(false)
// }

// ensureLocalLoad ensure local data file loaded
func (m *ToolManager) ensureLocalLoad(must bool) error {
	if m.localLoad {
		return nil
	}
	m.localLoad = true

	err := m.LoadLocalIndexes()
	if err != nil && must {
		panic(err)
	}
	return err
}

// LoadLocalIndexes local installed SDK,tool index information
func (m *ToolManager) LoadLocalIndexes() error {
	m.localFile = fsutil.ExpandHome("~/.xenv/tools.local.json")
	fileExist := fsutil.IsFile(m.localFile)
	xenvcom.Debugf("Load local index file: %s(exist=%v)\n", m.localFile, fileExist)

	if fileExist {
		err := jsonutil.DecodeFile(m.localFile, m.localTools)
		if err != nil {
			return err
		}
	}
	return nil
}

// FindLocalSdk find local installed sdk tool by name and version
func (m *ToolManager) FindLocalSdk(name, version string) *models.InstalledTool {
	_ = m.ensureLocalLoad(true)

	for _, tool := range m.localTools.SDKs {
		if tool.Name == name && tool.Version == version {
			return &tool
		}
	}
	return nil
}

// IndexLocalTools index local installed tools to datafile
func (m *ToolManager) IndexLocalTools() error {
	if err := m.ensureLocalLoad(false); err != nil {
		return err
	}

	currentTime := time.Now()
	if m.localTools.CreatedAt.IsZero() {
		m.localTools.CreatedAt = currentTime
	}
	m.localTools.UpdatedAt = currentTime
	m.localTools.SDKs = nil // 重新添加

	// SDK tools
	for _, sdkCfg := range m.config.SDKs {
		ccolor.Cyanf("Starting find installed %q SDK\n", sdkCfg.Name)

		if sdkCfg.InstallDir != "" {
			ver2dirMap, err := xenvutil.ListVersionDirs(sdkCfg.InstallDir)
			if err != nil {
				return err
			}

			baseDir := filepath.Dir(sdkCfg.InstallDir)
			ccolor.Cyanf(" - from dir: %s\n", baseDir)
			for version, installPath := range ver2dirMap {
				ccolor.Infof("  Found %s %s\n", sdkCfg.Name, version)

				// build local installed tool info
				m.localTools.SDKs = append(m.localTools.SDKs, models.InstalledTool{
					ID:     fmt.Sprintf("%s:%s", sdkCfg.Name, version),
					Name:   sdkCfg.Name,
					IsSDK:  true,
					BinDir: sdkCfg.BinDir,
					// version, install path
					Version:    version,
					InstallDir: installPath,
					CreatedAt:  currentTime,
				})
			}
		}

		// 不在统一目录下 InstallDir 的版本
		if sdkCfg.OtherVersions != nil {
			for version, dirPath := range sdkCfg.OtherVersions {
				dirPath = fsutil.ExpandHome(dirPath)
				if !fsutil.IsDir(dirPath) {
					ccolor.Warnf("[W] Custum version %s path %q is not exists\n", version, dirPath)
					continue
				}

				ccolor.Infof("  Found %s %s(at %s)\n", sdkCfg.Name, version, dirPath)
				// 添加到本地索引
				m.localTools.SDKs = append(m.localTools.SDKs, models.InstalledTool{
					ID:    fmt.Sprintf("%s:%s", sdkCfg.Name, version),
					Name:  sdkCfg.Name,
					IsSDK: true,
					// version, install path
					Version:    version,
					InstallDir: dirPath,
					CreatedAt:  currentTime,
				})
			}
		}
	}

	// TODO Simple tools

	ccolor.Magentaf("\nWrite indexed data to %s\n", m.localFile)
	return m.SaveLocalTools()
}

// SaveLocalTools saves the local tools information
func (m *ToolManager) SaveLocalTools() error {
	jsonBytes, err := json.MarshalIndent(m.localTools, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(m.localFile, jsonBytes, 0664)
}

func (m *ToolManager) AddSDKTool(name, version, installDir string) error {
	if err := m.ensureLocalLoad(false); err != nil {
		return err
	}

	// build local installed tool info
	currentTime := time.Now()
	if m.localTools.CreatedAt.IsZero() {
		m.localTools.CreatedAt = currentTime
	}
	m.localTools.UpdatedAt = currentTime

	m.localTools.SDKs = append(m.localTools.SDKs, models.InstalledTool{
		ID:         fmt.Sprintf("%s:%s", name, version),
		Name:       name,
		Version:    version,
		InstallDir: installDir,
		CreatedAt:  currentTime,
		UpdatedAt:  currentTime,
	})

	return m.SaveLocalTools()
}

func (m *ToolManager) DeleteSDKTool(localTool *models.InstalledTool) error {
	if err := m.ensureLocalLoad(false); err != nil {
		return err
	}

	// remove from ts.localTools
	sdkTools := m.localTools.SDKs
	toolIndex := localTool.Index
	m.localTools.SDKs = append(sdkTools[:toolIndex], sdkTools[toolIndex+1:]...)

	// save local.json
	return m.SaveLocalTools()
}

// FindSdkByID find local installed sdk tool by id
func (m *ToolManager) FindSdkByID(id string) *models.InstalledTool {
	_ = m.ensureLocalLoad(true)
	return m.localTools.FindSdkByID(id)
}

// ListSDKVersions 根据名称列出本地安装的SDK版本列表
func (m *ToolManager) ListSDKVersions(name string) []models.InstalledTool {
	// check caches
	if ls, ok := m.groupSdks[name]; ok {
		return ls
	}

	_ = m.ensureLocalLoad(true)
	ls := m.localTools.ListSdkByName(name)

	// cache for the sdk name
	if len(ls) > 0 {
		m.groupSdks[name] = ls
	}
	return ls
}

// MatchSdkByNameAndVersion 根据名称和版本匹配本地可用的一个sdk.
//
// 快捷方法，合并了 ListSDKVersions 和 MatchSdkByVersion
func (m *ToolManager) MatchSdkByNameAndVersion(name, version string) *models.InstalledTool {
	list := m.ListSDKVersions(name)
	if len(list) == 0 {
		return nil
	}
	return m.MatchSdkByVersion(list, version)
}

// MatchSdkByVersion 根据版本匹配本地可用的一个sdk
//
// 规则和优先级：
//   - 先完全匹配版本
//   - latest 匹配最新版本
//   - 1 可以匹配 1.19.x
//   - 1.19 可以匹配 1.19.x
//   - 1.19.x 可以匹配 1.19.x.x
func (m *ToolManager) MatchSdkByVersion(localSdks []models.InstalledTool, version string) *models.InstalledTool {
	dotNum := strings.Count(version, ".")

	// 完全匹配版本
	if dotNum > 1 {
		for _, localSdk := range localSdks {
			if localSdk.Version == version {
				return &localSdk
			}
		}
	}

	// latest 匹配最新版本 - 已排序，返回第一个
	if version == "latest" {
		return &localSdks[0]
	}

	// 前缀匹配
	for _, localSdk := range localSdks {
		locVersion := localSdk.Version
		// 检查版本是否以指定前缀开始，并且下一个字符是 '.' 或者字符串结束
		if strings.HasPrefix(locVersion, version) {
			if len(locVersion) == len(version) || locVersion[len(version)] == '.' {
				return &localSdk
			}
		}
	}

	// 设置了完整版本号，是否允许向上匹配版本
	if dotNum > 1 && m.config.AllowUpMatch > 0 {
		parts := strings.Split(version, ".")

		// 允许去除最后一位匹配 eg: 1.19.1 可以匹配 1.19.x
		if m.config.AllowUpMatch == xenvcom.UpMatchOne {
			matchVer := strings.Join(parts[:len(parts)-1], ".") + "."
			for _, localSdk := range localSdks {
				locVersion := localSdk.Version
				if strings.HasPrefix(locVersion, matchVer) {
					return &localSdk
				}
			}
		}
	}

	return nil
}

// LocalIndexes returns the local indexes
func (m *ToolManager) LocalIndexes() *models.ToolsLocal {
	return m.localTools
}
