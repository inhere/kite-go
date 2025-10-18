package manager

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/jsonutil"
	"github.com/gookit/goutil/x/ccolor"
	"github.com/inhere/kite-go/pkg/xenv/models"
	"github.com/inhere/kite-go/pkg/xenv/tools"
)

type ToolManager struct {
	init bool
	config *models.Configuration
	// local data file
	localLoad bool
	localFile  string
	localTools *models.ToolsLocal
	// tools register data - 从 config 配置中初始化 tools/config.json TODO
	configFile string
	configLoad bool
}

// NewToolManager creates a new ToolManager instance
func NewToolManager() *ToolManager {
	return &ToolManager{
		localTools: &models.ToolsLocal{Version: "v1"},
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

// ensureLocalLoad ensure local data file loaded
func (m *ToolManager) ensureLocalLoad(must bool) error {
	if m.localLoad {
		return nil
	}
	m.localLoad = true

	err := m.LoadLocalTools()
	if err != nil && must {
		panic(err)
	}
	return err
}

// LoadLocalTools local installed tools information
func (m *ToolManager) LoadLocalTools() error {
	m.localFile = fsutil.Expand("~/.xenv/tools.local.json")

	if fsutil.IsFile(m.localFile) {
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
	for _, toolCfg := range m.config.Tools {
		ver2dirMap, err := tools.ListVersionDirs(toolCfg.InstallDir)
		if err != nil {
			return err
		}

		ccolor.Cyanf("Found installed %q from %s\n", toolCfg.Name, toolCfg.InstallDir)
		for version, installPath := range ver2dirMap {
			ccolor.Infof("  Found %s %s\n", toolCfg.Name, version)
			// build local installed tool info
			m.localTools.SDKs = append(m.localTools.SDKs, models.InstalledTool{
				ID:         fmt.Sprintf("%s:%s", toolCfg.Name, version),
				Name:       toolCfg.Name,
				Version:    version,
				InstallDir: installPath,
				BinPaths:   []string{},
				CreatedAt:  currentTime,
			})
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
		BinPaths:   []string{},
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

func (m *ToolManager) FindSdkByID(id string) *models.InstalledTool {
	_ = m.ensureLocalLoad(true)
	return m.localTools.FindSdkByID(id)
}

func (m *ToolManager) ListLocalVersions(name string) []models.InstalledTool {
	_ = m.ensureLocalLoad(true)
	return m.localTools.ListSdkByName(name)
}
