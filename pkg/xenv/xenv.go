package xenv

import (
	"fmt"

	"github.com/inhere/kite-go/pkg/xenv/config"
	"github.com/inhere/kite-go/pkg/xenv/manager"
	"github.com/inhere/kite-go/pkg/xenv/service"
)

// ScriptMark 输出的脚本必须添加标记，前面部分为message, 后面部分为脚本
const ScriptMark = "--Expression--"

// Init initializes the xenv config and state
func Init() error {
	// Initialize configuration
	if err := config.Mgr.Init(); err != nil {
		return fmt.Errorf("failed to initialize configuration: %w", err)
	}

	if err := InitState(); err != nil {
		return fmt.Errorf("failed to initialize state manager: %w", err)
	}

	if err := toolMgr.Init(config.Config()); err != nil {
		return fmt.Errorf("failed to initialize tool manager: %w", err)
	}
	return nil
}

var stateMgr = manager.NewStateManager()

// State returns the state manager
func State() *manager.StateManager {
	return stateMgr
}

// InitState initializes the state manager
func InitState() error {
	return stateMgr.Init()
}

var toolMgr = manager.NewToolManager()

func ToolMgr() *manager.ToolManager {
	return toolMgr
}

func EnvService() (*service.EnvService, error) {
	// Initialize configuration
	if err := config.Mgr.Init(); err != nil {
		return nil, fmt.Errorf("failed to initialize configuration: %w", err)
	}

	if err := InitState(); err != nil {
		return nil, fmt.Errorf("failed to initialize state manager: %w", err)
	}

	// Create env manager
	return service.NewEnvService(config.Mgr.Config, stateMgr), nil
}

func ToolService() (*service.ToolService, error) {
	// Initialize configuration
	if err := Init(); err != nil {
		return nil, err
	}

	// Create tool service
	toolSvc := service.NewToolService(config.Config(), stateMgr, toolMgr)
	return toolSvc, nil
}
