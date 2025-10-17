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
	// DefaultShellHooksDir is the default directory for shell scripts
	DefaultShellHooksDir = "~/.config/xenv/hooks/"
)

// ConfigManager handles loading and saving configuration
type ConfigManager struct {
	Config *models.Configuration
}

// NewConfigManager creates a new ConfigManager with default configuration
func NewConfigManager() *ConfigManager {
	return &ConfigManager{
		Config: &models.Configuration{
			BinDir:          DefaultBinDir,
			InstallDir:      DefaultInstallDir,
			ShellHooksDir: DefaultShellHooksDir,
			Tools:           []models.ToolChain{},
			GlobalEnv: make(map[string]models.EnvVariable),
			GlobalPaths:     []models.PathEntry{},
		},
	}
}

// LoadConfig loads configuration from the specified file
func (cm *ConfigManager) LoadConfig(configPath string) error {
	cfg := config.New("xenv", config.WithTagName("json"))
	cfg.AddDriver(yaml.Driver)
	cfg.AddDriver(toml.Driver)

	// Load the configuration file
	err := cfg.LoadFiles(configPath)
	if err != nil {
		return err
	}

	// Load other configuration values like tools, global environment, etc.
	return cfg.Decode(&cm.Config)
}

// SaveConfig saves the configuration to the specified file
func (cm *ConfigManager) SaveConfig(configPath string) error {
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
