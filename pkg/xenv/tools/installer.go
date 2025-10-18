package tools

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gookit/goutil/strutil"
	"github.com/inhere/kite-go/internal/util"
	"github.com/inhere/kite-go/pkg/xenv/models"
)

// Installer handles the actual installation of tools
type Installer struct {
	config     *models.Configuration
	InstallDir string
}

// NewInstaller creates a new Installer
func NewInstaller(config *models.Configuration) *Installer {
	return &Installer{
		config: config,
	}
}

// EnsureBinDir ensures the bin directory exists and is in the PATH
func (i *Installer) EnsureBinDir() error {
	binDir := util.ExpandHome(i.config.BinDir)
	return util.EnsureDir(binDir)
}

// Install downloads and installs a tool
// 安装下载方式：http, git OR asdf, scoop, vfox, mise 等
func (i *Installer) Install(toolChain *models.ToolChain, version string) error {
	// Ensure bin directory exists
	if err := i.EnsureBinDir(); err != nil {
		return fmt.Errorf("failed to ensure bin directory: %w", err)
	}

	// If InstallURL is provided, download the tool
	if toolChain.InstallURL == "" {
		return fmt.Errorf("tool %q install_url is not configed", toolChain.Name)
	}
	if err := i.downloadAndExtract(toolChain, version); err != nil {
		return fmt.Errorf("failed to download and extract tool: %w", err)
	}

	// Create symlinks for executables
	return i.createShims(toolChain)
}

// downloadAndExtract downloads and extracts a tool based on its InstallURL
func (i *Installer) downloadAndExtract(toolChain *models.ToolChain, version string) (err error) {
	// Format the URL with version, os, and arch variables
	url := i.formatURL(toolChain.InstallURL, version)
	downExt := i.config.DownloadExt[runtime.GOOS]

	var tmpFile *os.File
	tmpDir := i.config.DownloadDir
	if tmpDir == "" {
		// Create temporary file for download
		tmpFile, err = os.CreateTemp("", fmt.Sprintf("xenv_%s_*", toolChain.Name))
		if err != nil {
			return fmt.Errorf("failed to create temp file: %w", err)
		}
		defer os.Remove(tmpFile.Name()) // Clean up
	} else {
		// Use the provided directory
		tmpFile, err = os.Open(tmpDir + "/" + toolChain.Name + "_" + version + "." + downExt)
		if err != nil {
			return fmt.Errorf("failed to open temp directory: %w", err)
		}
	}
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

	i.InstallDir = i.formatURL(toolChain.InstallDir, version)
	// Extract the file based on extension
	return i.extractFile(tmpFile.Name(), i.InstallDir)
}

// formatURL formats the installation URL with the appropriate variables
func (i *Installer) formatURL(urlOrPath, version string) string {
	downExt := i.config.DownloadExt[runtime.GOOS]
	varMap := map[string]string{
		"{os}":           runtime.GOOS,
		"{arch}":         runtime.GOARCH,
		"{version}":      version,
		"{download_ext}": downExt,
	}

	return strutil.Replaces(urlOrPath, varMap)
}

// extractFile extracts an archive file to the destination directory
func (i *Installer) extractFile(archivePath, destDir string) error {
	// Determine the file type based on extension
	if strings.HasSuffix(archivePath, ".zip") {
		return i.extractZip(archivePath, destDir)
	}

	if strings.HasSuffix(archivePath, ".tar.gz") || strings.HasSuffix(archivePath, ".tgz") {
		return i.extractTarGz(archivePath, destDir)
	}

	// If it's not an archive, just copy the file
	return i.copyExecutable(archivePath, destDir)
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
	binDir := util.ExpandHome(i.config.BinDir)

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
