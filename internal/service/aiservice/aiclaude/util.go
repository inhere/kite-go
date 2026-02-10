package aiclaude

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/sysutil"
)

var userConfigFile = ""

// UserConfigFile returns the path to the Claude user config file
//   - Linux, Mac: ~/.claude/config.json
//   - Windows: %USERPROFILE%\.claude\settings.json
func UserConfigFile() string {
	if userConfigFile == "" {
		filename := "config.json"
		if sysutil.IsWin() {
			filename = "settings.json"
		}
		userConfigFile = filepath.Join(fsutil.HomeDir(), ".claude", filename)
	}
	return userConfigFile
}

// UserClaudePath returns the path to the user's Claude directory
func UserClaudePath(subPath ...string) string {
	return fsutil.JoinPaths3(fsutil.HomeDir(), ".claude", subPath...)
}

// ReadUserConfig reads the Claude configuration from ~/.claude/config.json
func ReadUserConfig() (*ClaudeRuntimeConfig, error) {
	configPath := UserConfigFile()

	// If file doesn't exist, return empty config
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return &ClaudeRuntimeConfig{Env: make(map[string]string)}, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config ClaudeRuntimeConfig
	if err = json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	if config.Env == nil {
		config.Env = make(map[string]string)
	}
	return &config, nil
}

// WriteUserConfig writes the Claude configuration to ~/.claude/config.json with backup
func WriteUserConfig(config *ClaudeRuntimeConfig) error {
	configPath := UserConfigFile()

	// Ensure the directory exists
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Backup existing file if it exists
	if _, err := os.Stat(configPath); err == nil {
		backupPath := configPath + ".bak"
		if err := os.Rename(configPath, backupPath); err != nil {
			return fmt.Errorf("failed to backup config file: %w", err)
		}
	}

	// Write the config file
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err = os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	return nil
}
