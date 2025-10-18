package config

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Importer handles configuration import functionality
type Importer struct {
	configManager *Manager
}

// NewImporter creates a new Importer
func NewImporter(configManager *Manager) *Importer {
	return &Importer{
		configManager: configManager,
	}
}

// Import imports configuration from a file
func (i *Importer) Import(importPath string) error {
	// Check file extension to determine format
	ext := filepath.Ext(importPath)

	switch ext {
	case ".zip":
		return i.importFromZip(importPath)
	case ".json":
		return i.importFromJSON(importPath)
	default:
		// Try to detect format by reading the file
		return i.importFromFile(importPath)
	}
}

// importFromZip imports configuration from a ZIP file
func (i *Importer) importFromZip(importPath string) error {
	// Open the ZIP file
	zipReader, err := zip.OpenReader(importPath)
	if err != nil {
		return fmt.Errorf("failed to open ZIP file: %w", err)
	}
	defer zipReader.Close()

	// Find and read the config.json file from the ZIP
	for _, file := range zipReader.File {
		if filepath.Clean(file.Name) == "config.json" {
			// Open the config file inside the ZIP
			configFile, err := file.Open()
			if err != nil {
				return fmt.Errorf("failed to open config file in ZIP: %w", err)
			}
			defer configFile.Close()

			// Read the config data
			configData, err := io.ReadAll(configFile)
			if err != nil {
				return fmt.Errorf("failed to read config data: %w", err)
			}

			// Unmarshal the configuration
			var importedConfig struct {
				BinDir          string                                   `json:"bin_dir"`
				InstallDir      string                                   `json:"install_dir"`
				ShellHooksDir string `json:"shell_hooks_dir"`
				Tools           []interface{}                           `json:"tools"` // Will be handled separately
				GlobalEnv       map[string]map[string]interface{}       `json:"global_env"`
				GlobalPaths     []interface{}                           `json:"global_paths"`
			}

			if err := json.Unmarshal(configData, &importedConfig); err != nil {
				return fmt.Errorf("failed to unmarshal configuration: %w", err)
			}

			// Update the current configuration with imported values
			i.configManager.Config.BinDir = importedConfig.BinDir
			i.configManager.Config.InstallDir = importedConfig.InstallDir
			i.configManager.Config.ShellHooksDir = importedConfig.ShellHooksDir

			// For now, we'll store the raw JSON data and handle the conversion later
			// In a real implementation, we would convert the interface{} values to proper structs

			return nil
		}
	}

	return fmt.Errorf("config.json not found in ZIP file")
}

// importFromJSON imports configuration from a JSON file
func (i *Importer) importFromJSON(importPath string) error {
	// Read the file
	configData, err := os.ReadFile(importPath)
	if err != nil {
		return fmt.Errorf("failed to read configuration file: %w", err)
	}

	// Unmarshal the configuration
	var importedConfig struct {
		BinDir          string                                   `json:"bin_dir"`
		InstallDir      string                                   `json:"install_dir"`
		ShellHooksDir string `json:"shell_hooks_dir"`
		Tools           []interface{}                           `json:"tools"` // Will be handled separately
		GlobalEnv       map[string]map[string]interface{}       `json:"global_env"`
		GlobalPaths     []interface{}                           `json:"global_paths"`
	}

	if err := json.Unmarshal(configData, &importedConfig); err != nil {
		return fmt.Errorf("failed to unmarshal configuration: %w", err)
	}

	// Update the current configuration with imported values
	i.configManager.Config.BinDir = importedConfig.BinDir
	i.configManager.Config.InstallDir = importedConfig.InstallDir
	i.configManager.Config.ShellHooksDir = importedConfig.ShellHooksDir

	// For now, we'll store the raw JSON data and handle the conversion later
	// In a real implementation, we would convert the interface{} values to proper structs

	return nil
}

// importFromFile attempts to import from any file by detecting its format
func (i *Importer) importFromFile(importPath string) error {
	// Try to detect if it's a JSON file by reading the beginning
	file, err := os.Open(importPath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Read the first few bytes to check if it's JSON
	buf := make([]byte, 2)
	_, err = file.Read(buf)
	if err != nil && err != io.EOF {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Reset file pointer
	_, err = file.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("failed to reset file pointer: %w", err)
	}

	// Check if it starts with '{', which would indicate JSON
	if len(buf) >= 1 && buf[0] == '{' {
		return i.importFromJSON(importPath)
	}

	// If not clearly JSON, assume it's a ZIP if it has the ZIP signature
	// ZIP files start with 0x50, 0x4B
	if len(buf) >= 2 && buf[0] == 0x50 && buf[1] == 0x4B {
		return i.importFromZip(importPath)
	}

	return fmt.Errorf("unable to detect file format for %s", importPath)
}
