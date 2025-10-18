package config

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"os"
)

// Exporter handles configuration export functionality
type Exporter struct {
	configManager *Manager
}

// NewExporter creates a new Exporter
func NewExporter(configManager *Manager) *Exporter {
	return &Exporter{
		configManager: configManager,
	}
}

// Export exports the configuration to a file in the specified format
func (e *Exporter) Export(exportPath string, format string) error {
	switch format {
	case "zip":
		return e.exportToZip(exportPath)
	case "json":
		return e.exportToJSON(exportPath)
	default:
		return fmt.Errorf("unsupported export format: %s", format)
	}
}

// exportToZip exports configuration to a ZIP file
func (e *Exporter) exportToZip(exportPath string) error {
	// Create a new zip file
	zipFile, err := os.Create(exportPath)
	if err != nil {
		return fmt.Errorf("failed to create export file: %w", err)
	}
	defer zipFile.Close()

	// Create a new zip writer
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Add the main configuration to the ZIP
	configData, err := json.MarshalIndent(e.configManager.Config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal configuration: %w", err)
	}

	configFile, err := zipWriter.Create("config.json")
	if err != nil {
		return fmt.Errorf("failed to create config entry in zip: %w", err)
	}
	if _, err := configFile.Write(configData); err != nil {
		return fmt.Errorf("failed to write config data to zip: %w", err)
	}

	// Check file size
	fileInfo, err := zipFile.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	// Verify size limit (10MB)
	const maxSize = 10 * 1024 * 1024 // 10MB in bytes
	if fileInfo.Size() > maxSize {
		return fmt.Errorf("exported file exceeds size limit of 10MB: %d bytes", fileInfo.Size())
	}

	return nil
}

// exportToJSON exports configuration to a JSON file
func (e *Exporter) exportToJSON(exportPath string) error {
	// Marshal the configuration to JSON
	configData, err := json.MarshalIndent(e.configManager.Config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal configuration: %w", err)
	}

	// Write the JSON to the file
	if err := os.WriteFile(exportPath, configData, 0644); err != nil {
		return fmt.Errorf("failed to write configuration to file: %w", err)
	}

	// Verify size limit (10MB)
	fileInfo, err := os.Stat(exportPath)
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	const maxSize = 10 * 1024 * 1024 // 10MB in bytes
	if fileInfo.Size() > maxSize {
		return fmt.Errorf("exported file exceeds size limit of 10MB: %d bytes", fileInfo.Size())
	}

	return nil
}

// ExportWithRelatedFiles exports configuration and related files to a ZIP
func (e *Exporter) ExportWithRelatedFiles(exportPath string, includeRelated bool) error {
	if !includeRelated {
		return e.exportToZip(exportPath)
	}

	// Create a new zip file
	zipFile, err := os.Create(exportPath)
	if err != nil {
		return fmt.Errorf("failed to create export file: %w", err)
	}
	defer zipFile.Close()

	// Create a new zip writer
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Add the main configuration to the ZIP
	configData, err := json.MarshalIndent(e.configManager.Config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal configuration: %w", err)
	}

	configFile, err := zipWriter.Create("config.json")
	if err != nil {
		return fmt.Errorf("failed to create config entry in zip: %w", err)
	}
	if _, err := configFile.Write(configData); err != nil {
		return fmt.Errorf("failed to write config data to zip: %w", err)
	}

	// Add related configuration files (like activity state)
	activityState, err := loadActivityState()
	if err == nil {
		activityData, err := json.MarshalIndent(activityState, "", "  ")
		if err == nil {
			activityFile, err := zipWriter.Create("activity.json")
			if err == nil {
				activityFile.Write(activityData)
			}
		}
	}

	// Check file size
	fileInfo, err := zipFile.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	// Verify size limit (10MB)
	const maxSize = 10 * 1024 * 1024 // 10MB in bytes
	if fileInfo.Size() > maxSize {
		return fmt.Errorf("exported file exceeds size limit of 10MB: %d bytes", fileInfo.Size())
	}

	return nil
}

// loadActivityState loads the activity state (stub implementation - would need to be imported from models)
func loadActivityState() (interface{}, error) {
	// This is a placeholder implementation
	// In a real implementation, we would need to import and use the models package
	return nil, fmt.Errorf("not implemented")
}
