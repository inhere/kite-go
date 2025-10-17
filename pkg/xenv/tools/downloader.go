package tools

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// Downloader handles downloading tools from external sources
type Downloader struct {
}

// NewDownloader creates a new Downloader
func NewDownloader() *Downloader {
	return &Downloader{}
}

// DownloadFile downloads a file from the given URL to the specified destination
func (d *Downloader) DownloadFile(url, destPath string) error {
	// Ensure the destination directory exists
	destDir := filepath.Dir(destPath)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Create the destination file
	destFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destFile.Close()

	return d.doDownload(url, destFile)
}

// DownloadToTemp downloads a file to a temporary location
func (d *Downloader) DownloadToTemp(url, prefix string) (string, error) {
	// Create temporary file
	tmpFile, err := os.CreateTemp("", prefix)
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()
	defer tmpFile.Close() // Close but keep the file

	err = d.doDownload(url, tmpFile)
	if err != nil {
		_ = os.Remove(tmpPath) // Clean up on error
	}
	return tmpPath, err
}

// doDownload downloads a file from the given URL to the specified destination
func (d *Downloader) doDownload(url string, destFile *os.File) error {
	// Download the file
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status: %d", resp.StatusCode)
	}

	// Copy response body to destination file
	_, err = io.Copy(destFile, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save downloaded file: %w", err)
	}

	return nil
}
