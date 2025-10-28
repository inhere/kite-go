package manager

import (
	"fmt"
	"os"

	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/strutil"
	"github.com/inhere/kite-go/pkg/util"
	"github.com/inhere/kite-go/pkg/xenv/models"
)

// DirenvManager handles directory-level configuration files like .xenv.toml or .envrc
type DirenvManager struct {
}

// EditTomlFile edits a TOML file
func (m *DirenvManager) EditTomlFile(state *models.ActivityState) error {

}

// ProcessDirectoryConfig processes directory-level configuration files like .xenv.toml or .envrc
func ProcessDirectoryConfig() error {
	// Get the current working directory
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %w", err)
	}

	// Check for .xenv.toml file in the current directory and parent directories up to the root
	xenvTomlPath := fsutil.FindOneInParentDirs(wd, ".xenv.toml")
	if xenvTomlPath != "" {
		fmt.Printf("Found .xenv.toml at: %s\n", xenvTomlPath)
		// Process the .xenv.toml file
		if err := processXenvToml(xenvTomlPath); err != nil {
			return fmt.Errorf("failed to process .xenv.toml: %w", err)
		}
	}

	// Check for .envrc file in the current directory and parent directories up to the root
	envrcPath := fsutil.FindOneInParentDirs(wd, strutil.OrCond(util.IsHookBash(), ".envrc", ".envrc.ps1"))
	if envrcPath != "" {
		fmt.Printf("Found .envrc at: %s\n", envrcPath)
		// Process the .envrc file
		if err := processEnvrc(envrcPath); err != nil {
			return fmt.Errorf("failed to process .envrc: %w", err)
		}
	}

	return nil
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
