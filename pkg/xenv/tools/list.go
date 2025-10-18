package tools

import (
	"fmt"

	"github.com/inhere/kite-go/pkg/xenv/models"
)

// List lists installed tools
type List struct {
	service *ToolService
}

// NewList creates a new List
func NewList(service *ToolService) *List {
	return &List{
		service: service,
	}
}

// ListAll lists all tools
func (l *List) ListAll(showAll bool) {
	cfgTools := l.service.ListTools()
	if len(cfgTools) == 0 {
		fmt.Println("No tools for managed. see config: tools, simple_tools")
		return
	}

	fmt.Println("Installed tools:")

	for _, tool := range cfgTools {
		status := ""
		if tool.Installed {
			status = " [INSTALLED]"
		} else {
			status = " [NOT INSTALLED]"
		}

		fmt.Printf("- %s %s\n", tool.Name, status)
		if showAll {
			fmt.Printf("  InstallDir: %s\n", tool.InstallDir)
			fmt.Printf("  BinPaths: %v\n", tool.BinPaths)
			if len(tool.Alias) > 0 {
				fmt.Printf("  Aliases: %v\n", tool.Alias)
			}
		}
	}
}

// ListTools returns tools based on filters
func (l *List) ListTools(filterType string) []models.ToolChain {
	allTools := l.service.ListTools()

	switch filterType {
	case "tool", "":
		// Return all installed tools (default behavior)
		return allTools
	case "env":
		// This would return environment-related items, but that's not a tool
		// This functionality might be handled by env module
		return []models.ToolChain{}
	case "path":
		// This would return PATH-related items, but that's not a tool
		// This functionality might be handled by env module
		return []models.ToolChain{}
	case "activity":
		// This would return active tools - would need to get from ActivityState
		return allTools // Placeholder - would need to filter based on active state
	default:
		return allTools
	}
}
