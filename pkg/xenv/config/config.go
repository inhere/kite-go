package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/toml"
	"github.com/gookit/config/v2/yaml"
	"github.com/inhere/kite-go/pkg/xenv/models"
)

const (
	// DefaultBinDir is the default directory for installed tools
	DefaultBinDir = "~/.local/bin"
	// DefaultInstallDir is the default directory for tool installation
	DefaultInstallDir = "~/.xenv/tools"
	// DefaultConfigDir is the default config directory for xenv
	DefaultConfigDir = "~/.config/xenv/"
	// DefaultShellHooksDir is the default directory for shell scripts
	DefaultShellHooksDir = "~/.config/xenv/hooks/"
)

// Manager handles loading and saving configuration
type Manager struct {
	cfgInit bool
	Config *models.Configuration
}

// Mgr is the global ConfigManager instance
var Mgr = NewConfigManager()

func Config() *models.Configuration {
	return Mgr.Config
}

// NewConfigManager creates a new ConfigManager with default configuration
func NewConfigManager() *Manager {
	return &Manager{
		Config: &models.Configuration{
			BinDir:       DefaultBinDir,
			InstallDir:   DefaultInstallDir,
			ShellHooksDir: DefaultShellHooksDir,
			ShellAliases: make(map[string]string),
			// env
			GlobalEnv:   make(map[string]string),
			GlobalPaths: []string{},
			// tools
			Tools:       []models.ToolChain{},
			SimpleTools: []models.SimpleTool{},
			DownloadDir: DefaultInstallDir + "/cache",
			DownloadExt: map[string]string{
				"windows": "zip",
				"linux":   "tar.gz",
				"darwin":  "tar.gz",
			},
		},
	}
}

// Init initializes load the configuration data
func (cm *Manager) Init() error {
	if cm.cfgInit {
		return nil
	}
	cm.cfgInit = true
	return cm.LoadConfig(GetDefaultConfigPath())
}

// LoadConfig loads configuration from the specified file
func (cm *Manager) LoadConfig(configPath string) error {
	cfg := config.New("xenv", config.WithTagName("json"))
	cfg.AddDriver(yaml.Driver)
	cfg.AddDriver(toml.Driver)

	// Load the configuration file
	err := cfg.LoadFiles(configPath)
	if err != nil {
		return err
	}

	// Load other configuration values like tools, global environment, etc.
	err = cfg.Decode(&cm.Config)
	cm.Config.SetConfigFile(configPath)
	cm.Config.SetConfigDir(GetDefaultConfigDir())
	return err
}

// SaveConfig saves the configuration to the specified file
func (cm *Manager) SaveConfig(configPath string) error {
	// Ensure the config directory exists
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	// Marshal the configuration to JSON
	configData, err := json.MarshalIndent(cm.Config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal configuration to JSON: %w", err)
	}

	// Write to the file
	return os.WriteFile(configPath, configData, 0644)
}

// GetDefaultConfigPath returns the default configuration file path
func GetDefaultConfigPath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".config", "xenv", "config.yaml")
}

func GetDefaultConfigDir() string {
	// Get the config directory
	homeDir, _ := os.UserHomeDir()
	// if err != nil {
	// 	return fmt.Errorf("failed to get user home directory: %w", err)
	// }
	return filepath.Join(homeDir, ".config", "xenv")
}
