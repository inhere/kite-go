package aiclaude

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// ClaudeRuntimeConfig represents the Claude configuration file format
//
//  - Linux, Mac: ~/.claude/config.json
//  - Windows: %USERPROFILE%\.claude\settings.json
type ClaudeRuntimeConfig struct {
	Env map[string]string `json:"env,omitempty"`
	// IncludeCoAuthoredBy indicates whether to include co-authored-by in the commit message
	IncludeCoAuthoredBy bool `json:"includeCoAuthoredBy"`
	// enabledPlugins map
	EnabledPlugins map[string]any `json:"enabledPlugins,omitempty"`
	// statusLine string map
	StatusLine map[string]string `json:"statusLine,omitempty"`
	// configFile path of the configuration file(internal)
	configFile string
	fileExists bool
}

// ConfigFile returns the path of the configuration file
func (c *ClaudeRuntimeConfig) ConfigFile() string {
	return c.configFile
}

// Save saves the configuration to the file
func (c *ClaudeRuntimeConfig) Save() error {
	// Ensure the directory exists
	dir := filepath.Dir(c.configFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Backup existing file if it exists
	if _, err := os.Stat(c.configFile); err == nil {
		backupPath := c.configFile + ".bak"
		if err := os.Rename(c.configFile, backupPath); err != nil {
			return fmt.Errorf("failed to backup config file: %w", err)
		}
	}

	// Write the config file
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err = os.WriteFile(c.configFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	return nil
}
