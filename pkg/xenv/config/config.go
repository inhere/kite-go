package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gookit/config/v2"
	"github.com/inhere/kite-go/pkg/xenv/models"
)

const (
	// DefaultBinDir is the default directory for installed tools
	DefaultBinDir = "~/.local/bin"
	// DefaultInstallDir is the default directory for tool installation
	DefaultInstallDir = "~/.xenv/tools"
	// DefaultShellScriptsDir is the default directory for shell scripts
	DefaultShellScriptsDir = "~/.config/xenv/hooks/"
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
			ShellScriptsDir: DefaultShellScriptsDir,
			Tools:           []models.ToolChain{},
			GlobalEnv:       make(map[string]models.EnvironmentVariable),
			GlobalPaths:     []models.PathEntry{},
		},
	}
}

// LoadConfig loads configuration from the specified file
func (cm *ConfigManager) LoadConfig(configPath string) error {
	// Load the configuration file
	err := config.LoadFiles(configPath)
	if err != nil {
		return err
	}

	// Get the loaded config data
	loadedConfig := config.Data()

	// Map loaded config to our Configuration model
	if binDir, ok := loadedConfig["bin_dir"].(string); ok {
		cm.Config.BinDir = binDir
	}
	if installDir, ok := loadedConfig["install_dir"].(string); ok {
		cm.Config.InstallDir = installDir
	}
	if shellScriptsDir, ok := loadedConfig["shell_scripts_dir"].(string); ok {
		cm.Config.ShellScriptsDir = shellScriptsDir
	}

	// TODO: Load other configuration values like tools, global environment, etc.

	return nil
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