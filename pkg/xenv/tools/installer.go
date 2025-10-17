package tools

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/inhere/kite-go/internal/util"
	"github.com/inhere/kite-go/pkg/xenv/models"
)

// Installer handles the actual installation of tools
type Installer struct {
	service *ToolService
}

// NewInstaller creates a new Installer
func NewInstaller(service *ToolService) *Installer {
	return &Installer{
		service: service,
	}
}

// Install downloads and installs a tool
func (i *Installer) Install(toolChain *models.ToolChain) error {
	// Ensure bin directory exists
	if err := i.service.EnsureBinDir(); err != nil {
		return fmt.Errorf("failed to ensure bin directory: %w", err)
	}

	// If InstallURL is provided, download the tool
	if toolChain.InstallURL != "" {
		if err := i.downloadAndExtract(toolChain); err != nil {
			return fmt.Errorf("failed to download and extract tool: %w", err)
		}
	} else {
		// If no InstallURL, assume the tool is already available locally
		// This is for manually installed tools
		fmt.Printf("Tool %s marked as installed (local installation)\n", toolChain.ID)
	}

	// Create symlinks for executables
	return i.createShims(toolChain)
}

// downloadAndExtract downloads and extracts a tool based on its InstallURL
func (i *Installer) downloadAndExtract(toolChain *models.ToolChain) error {
	// Format the URL with version, os, and arch variables
	url := i.formatURL(toolChain.InstallURL, toolChain.Version)

	// Create temporary file for download
	tmpFile, err := os.CreateTemp("", fmt.Sprintf("xenv_%s_*", toolChain.Name))
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name()) // Clean up
	defer tmpFile.Close()

	// Download the file
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download tool: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status: %d", resp.StatusCode)
	}

	// Copy response body to temp file
	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save downloaded file: %w", err)
	}

	// Extract the file based on extension
	return i.extractFile(tmpFile.Name(), toolChain.InstallDir)
}

// formatURL formats the installation URL with the appropriate variables
func (i *Installer) formatURL(template, version string) string {
	result := template
	result = replaceAll(result, "{version}", version)
	result = replaceAll(result, "{os}", runtime.GOOS)
	result = replaceAll(result, "{arch}", runtime.GOARCH)
	
	// Simple Windows check for file extension
	if runtime.GOOS == "windows" {
		result = replaceAll(result, "{isWindows ? zip : tar.gz}", "zip")
		result = replaceAll(result, "{isWindows ? .exe : }", ".exe")
	} else {
		result = replaceAll(result, "{isWindows ? zip : tar.gz}", "tar.gz")
		result = replaceAll(result, "{isWindows ? .exe : }", "")
	}
	
	return result
}

// replaceAll replaces all occurrences of old with new in s
func replaceAll(s, old, new string) string {
	return filepath.ToSlash(filepath.Clean(strings.Replace(s, old, new, -1)))
}

// extractFile extracts an archive file to the destination directory
func (i *Installer) extractFile(archivePath, destDir string) error {
	// Determine the file type based on extension
	if strings.HasSuffix(archivePath, ".zip") {
		return i.extractZip(archivePath, destDir)
	} else if strings.HasSuffix(archivePath, ".tar.gz") || strings.HasSuffix(archivePath, ".tgz") {
		return i.extractTarGz(archivePath, destDir)
	} else {
		// If it's not an archive, just copy the file
		return i.copyExecutable(archivePath, destDir)
	}
}

// extractZip extracts a ZIP archive
func (i *Installer) extractZip(archivePath, destDir string) error {
	// TODO: Implement ZIP extraction
	// This is a placeholder for the actual implementation
	return nil
}

// extractTarGz extracts a tar.gz archive
func (i *Installer) extractTarGz(archivePath, destDir string) error {
	// TODO: Implement tar.gz extraction
	// This is a placeholder for the actual implementation
	return nil
}

// copyExecutable copies a single executable file to the destination
func (i *Installer) copyExecutable(srcPath, destDir string) error {
	// Get the file name
	_, fileName := filepath.Split(srcPath)
	destPath := filepath.Join(destDir, fileName)
	
	// Copy the file
	return util.CopyFile(srcPath, destPath)
}

// createShims creates symlinks (shims) for the tool executables
func (i *Installer) createShims(toolChain *models.ToolChain) error {
	binDir := util.ExpandHome(i.service.config.BinDir)
	
	// For each binary path of the tool, create a shim
	for _, binPath := range toolChain.BinPaths {
		// Get all executable files in the bin path
		entries, err := os.ReadDir(binPath)
		if err != nil {
			continue // Skip if directory doesn't exist
		}
		
		for _, entry := range entries {
			if !entry.IsDir() && isExecutable(entry.Name()) {
				// Create a symlink in the bin directory
				srcPath := filepath.Join(binPath, entry.Name())
				dstPath := filepath.Join(binDir, entry.Name())
				
				// Create the shim (symbolic link)
				if err := util.CreateSymlink(srcPath, dstPath); err != nil {
					// If symlinks aren't supported (e.g., on Windows without admin), copy instead
					if runtime.GOOS == "windows" {
						if copyErr := util.CopyFile(srcPath, dstPath); copyErr != nil {
							fmt.Printf("Warning: Failed to copy or link %s: %v\n", entry.Name(), err)
						}
					} else {
						fmt.Printf("Warning: Failed to create symlink for %s: %v\n", entry.Name(), err)
					}
				}
			}
		}
	}
	
	return nil
}

// isExecutable checks if a file name suggests it's executable
func isExecutable(filename string) bool {
	if runtime.GOOS == "windows" {
		return strings.HasSuffix(strings.ToLower(filename), ".exe") || 
			   strings.HasSuffix(strings.ToLower(filename), ".bat") || 
			   strings.HasSuffix(strings.ToLower(filename), ".cmd") ||
			   strings.HasSuffix(strings.ToLower(filename), ".com")
	}
	// For Unix-like systems, we can't determine executability just from the name
	// We'll assume files without extensions are potentially executable
	return !strings.Contains(filename, ".")
}