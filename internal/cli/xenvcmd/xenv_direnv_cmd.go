package xenvcmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gookit/gcli/v3"
)

// ProcessDirectoryConfig processes directory-level configuration files like .xenv.toml or .envrc
func ProcessDirectoryConfig() error {
	// Get the current working directory
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %w", err)
	}

	// Check for .xenv.toml file in the current directory and parent directories up to the root
	xenvTomlPath, err := findFileInParentDirs(wd, ".xenv.toml")
	if err == nil && xenvTomlPath != "" {
		fmt.Printf("Found .xenv.toml at: %s\n", xenvTomlPath)
		// Process the .xenv.toml file
		if err := processXenvToml(xenvTomlPath); err != nil {
			return fmt.Errorf("failed to process .xenv.toml: %w", err)
		}
	}

	// Check for .envrc file in the current directory and parent directories up to the root
	envrcPath, err := findFileInParentDirs(wd, ".envrc")
	if err == nil && envrcPath != "" {
		fmt.Printf("Found .envrc at: %s\n", envrcPath)
		// Process the .envrc file
		if err := processEnvrc(envrcPath); err != nil {
			return fmt.Errorf("failed to process .envrc: %w", err)
		}
	}

	return nil
}

// findFileInParentDirs looks for a file in the current directory and parent directories
func findFileInParentDirs(startDir, fileName string) (string, error) {
	currentDir := startDir

	for {
		// Check if the file exists in the current directory
		filePath := filepath.Join(currentDir, fileName)
		if _, err := os.Stat(filePath); err == nil {
			// File found
			return filePath, nil
		}

		// Get parent directory
		parentDir := filepath.Dir(currentDir)

		// If we reached the root directory, stop searching
		if parentDir == currentDir {
			// Reached the root, file not found
			return "", nil
		}

		// Move to parent directory
		currentDir = parentDir
	}
}

// processXenvToml processes an .xenv.toml file
func processXenvToml(filePath string) error {
	// This is a placeholder implementation
	// In a real implementation, we would parse the TOML file and apply the configuration
	fmt.Printf("Processing .xenv.toml file: %s (not implemented yet)\n", filePath)
	return nil
}

// processEnvrc processes an .envrc file
func processEnvrc(filePath string) error {
	// This is a placeholder implementation
	// In a real implementation, we would source the .envrc file
	fmt.Printf("Processing .envrc file: %s (not implemented yet)\n", filePath)
	return nil
}

// DirectoryConfigCmd implements a command to process directory configuration
var DirectoryConfigCmd = &gcli.Command{
	Name:    "direnv",
	Desc:   "Process directory-level configuration files",
	Hidden:  true, // This is more of an internal command
	Func: func(c *gcli.Command, args []string) error {
		return ProcessDirectoryConfig()
	},
}
