package models

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/gookit/goutil/fsutil"
)

// ActivityState 代表用户当前激活的工具链和环境状态
type ActivityState struct {
	ID           string            `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ActiveTools  map[string]string `json:"active_tools"`  // 激活的工具链映射，key为工具名，value为版本
	ActiveEnv    map[string]string `json:"active_env"`    // 激活的环境变量
	ActivePaths  []string          `json:"active_paths"`  // 激活的路径列表
}

func NewActivityState() *ActivityState {
	return &ActivityState{
		ID:          "default",
		ActiveTools: make(map[string]string),
		ActiveEnv:   make(map[string]string),
		ActivePaths: []string{},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// LoadGlobalState loads the activity state from file
func LoadGlobalState() (*ActivityState, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	statePath := filepath.Join(homeDir, ".config", "xenv", "activity.json")

	// Check if file exists
	if _, err := os.Stat(statePath); os.IsNotExist(err) {
		// Return default state if file doesn't exist
		return NewActivityState(), nil
	}

	// Read the file
	data, err := os.ReadFile(statePath)
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON
	var state ActivityState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, err
	}
	return &state, nil
}

// Save saves the activity state to file
func (as *ActivityState) Save(filePath string) error {
	if err := fsutil.MkParentDir(filePath); err != nil {
		return err
	}

	// Update timestamps
	as.UpdatedAt = time.Now()
	if as.CreatedAt.IsZero() {
		as.CreatedAt = as.UpdatedAt
	}

	// Marshal to JSON
	data, err := json.MarshalIndent(as, "", "  ")
	if err != nil {
		return err
	}

	// Write to file
	return os.WriteFile(filePath, data, 0644)
}
