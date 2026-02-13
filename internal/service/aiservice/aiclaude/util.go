package aiclaude

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gookit/goutil/fsutil"
)

var userConfigFile = ""

// UserConfigFile returns the path to the Claude user config file
//   - Linux, Mac: ~/.claude/config.json
//   - Windows: %USERPROFILE%\.claude\settings.json
func UserConfigFile() string {
	if userConfigFile == "" {
		userConfigFile = filepath.Join(fsutil.HomeDir(), ".claude", "settings.json")
	}
	return userConfigFile
}

// findConfigFile finds the path to the Claude configuration file
//
//  - user: ~/.claude/settings.json
//  - project:
//     - .claudue/settings.json
//     - .claudue/settings.local.json
func findConfigFile(scope string) string {
	if scope == "user" {
		return UserConfigFile()
	}

	// scope == "project"
	localFile := ".claude/settings.local.json"
	if fsutil.IsFile(localFile) {
		return localFile
	}
	return ".claude/settings.json"
}

// UserClaudePath returns the path to the user's Claude directory
func UserClaudePath(subPath ...string) string {
	return fsutil.JoinPaths3(fsutil.HomeDir(), ".claude", subPath...)
}

// LoadConfig reads the Claude configuration by scope
func LoadConfig(scope string) (*ClaudeRuntimeConfig, error) {
	configFile := findConfigFile(scope)

	// If file doesn't exist, return empty config
	data, err := os.ReadFile(configFile)
	if err != nil {
		return &ClaudeRuntimeConfig{
			Env:        make(map[string]string),
			configFile: configFile,
		}, nil
	}

	var config ClaudeRuntimeConfig
	if err = json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	config.fileExists = true
	config.configFile = configFile
	if config.Env == nil {
		config.Env = make(map[string]string)
	}
	return &config, nil
}

