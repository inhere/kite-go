package envmgr

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gookit/goutil/fsutil"
)

// DefaultStateManager 默认状态管理器实现
type DefaultStateManager struct {
	stateFile string
	mutex     sync.RWMutex
}

// NewStateManager 创建状态管理器
func NewStateManager(stateFile string) *DefaultStateManager {
	return &DefaultStateManager{
		stateFile: stateFile,
	}
}

// LoadState 加载活跃状态
func (sm *DefaultStateManager) LoadState() (*ActiveState, error) {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	if !fsutil.IsFile(sm.stateFile) {
		return sm.getDefaultState(), nil
	}

	data, err := os.ReadFile(sm.stateFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read state file %s: %w", sm.stateFile, err)
	}

	var state ActiveState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("failed to parse state file %s: %w", sm.stateFile, err)
	}

	// 确保字段不为nil
	if state.CurrentSDKs == nil {
		state.CurrentSDKs = make(map[string]string)
	}
	if state.AddPaths == nil {
		state.AddPaths = []string{}
	}
	if state.AddEnvs == nil {
		state.AddEnvs = make(map[string]string)
	}

	return &state, nil
}

// SaveState 保存活跃状态
func (sm *DefaultStateManager) SaveState(state *ActiveState) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	if state == nil {
		return fmt.Errorf("state is nil")
	}

	// 更新时间戳
	state.UpdatedAt = time.Now()

	// 确保目录存在
	if err := fsutil.MkdirQuick(filepath.Dir(sm.stateFile)); err != nil {
		return fmt.Errorf("failed to create state directory: %w", err)
	}

	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}

	// 写入临时文件，然后原子性重命名
	tempFile := sm.stateFile + ".tmp"
	if err := os.WriteFile(tempFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write temp state file: %w", err)
	}

	if err := os.Rename(tempFile, sm.stateFile); err != nil {
		os.Remove(tempFile) // 清理临时文件
		return fmt.Errorf("failed to rename temp state file: %w", err)
	}

	return nil
}

// UpdateSDKState 更新SDK状态
func (sm *DefaultStateManager) UpdateSDKState(sdk, version string, active bool) error {
	state, err := sm.LoadState()
	if err != nil {
		return err
	}

	if active {
		state.CurrentSDKs[sdk] = version
	} else {
		delete(state.CurrentSDKs, sdk)
	}

	return sm.SaveState(state)
}

// GetCurrentSDKs 获取当前激活的SDK
func (sm *DefaultStateManager) GetCurrentSDKs() (map[string]string, error) {
	state, err := sm.LoadState()
	if err != nil {
		return nil, err
	}

	// 返回副本避免外部修改
	result := make(map[string]string)
	for k, v := range state.CurrentSDKs {
		result[k] = v
	}

	return result, nil
}

// AddPath 添加路径到状态
func (sm *DefaultStateManager) AddPath(path string) error {
	state, err := sm.LoadState()
	if err != nil {
		return err
	}

	// 检查路径是否已存在
	for _, existingPath := range state.AddPaths {
		if existingPath == path {
			return nil // 已存在，无需添加
		}
	}

	state.AddPaths = append(state.AddPaths, path)
	return sm.SaveState(state)
}

// RemovePath 从状态中移除路径
func (sm *DefaultStateManager) RemovePath(path string) error {
	state, err := sm.LoadState()
	if err != nil {
		return err
	}

	for i, existingPath := range state.AddPaths {
		if existingPath == path {
			state.AddPaths = append(state.AddPaths[:i], state.AddPaths[i+1:]...)
			break
		}
	}

	return sm.SaveState(state)
}

// SetEnv 设置环境变量到状态
func (sm *DefaultStateManager) SetEnv(name, value string) error {
	state, err := sm.LoadState()
	if err != nil {
		return err
	}

	state.AddEnvs[name] = value
	return sm.SaveState(state)
}

// UnsetEnv 从状态中移除环境变量
func (sm *DefaultStateManager) UnsetEnv(name string) error {
	state, err := sm.LoadState()
	if err != nil {
		return err
	}

	delete(state.AddEnvs, name)
	return sm.SaveState(state)
}

// ClearState 清空状态
func (sm *DefaultStateManager) ClearState() error {
	state := sm.getDefaultState()
	return sm.SaveState(state)
}

// GetStateFilePath 获取状态文件路径
func (sm *DefaultStateManager) GetStateFilePath() string {
	return sm.stateFile
}

// IsStateFileExists 检查状态文件是否存在
func (sm *DefaultStateManager) IsStateFileExists() bool {
	return fsutil.IsFile(sm.stateFile)
}

// BackupState 备份当前状态
func (sm *DefaultStateManager) BackupState() error {
	if !sm.IsStateFileExists() {
		return fmt.Errorf("state file does not exist")
	}

	backupFile := sm.stateFile + ".backup." + time.Now().Format("20060102-150405")

	data, err := os.ReadFile(sm.stateFile)
	if err != nil {
		return fmt.Errorf("failed to read state file for backup: %w", err)
	}

	if err := os.WriteFile(backupFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write backup file: %w", err)
	}

	return nil
}

// RestoreState 从备份恢复状态
func (sm *DefaultStateManager) RestoreState(backupFile string) error {
	if !fsutil.IsFile(backupFile) {
		return fmt.Errorf("backup file does not exist: %s", backupFile)
	}

	data, err := os.ReadFile(backupFile)
	if err != nil {
		return fmt.Errorf("failed to read backup file: %w", err)
	}

	// 验证备份文件格式
	var state ActiveState
	if err := json.Unmarshal(data, &state); err != nil {
		return fmt.Errorf("invalid backup file format: %w", err)
	}

	if err := os.WriteFile(sm.stateFile, data, 0644); err != nil {
		return fmt.Errorf("failed to restore state file: %w", err)
	}

	return nil
}

// GetStateStats 获取状态统计信息
func (sm *DefaultStateManager) GetStateStats() (*StateStats, error) {
	state, err := sm.LoadState()
	if err != nil {
		return nil, err
	}

	stats := &StateStats{
		ActiveSDKCount: len(state.CurrentSDKs),
		PathCount:      len(state.AddPaths),
		EnvCount:       len(state.AddEnvs),
		LastUpdated:    state.UpdatedAt,
	}

	// 计算文件大小
	if fsutil.IsFile(sm.stateFile) {
		if info, err := os.Stat(sm.stateFile); err == nil {
			stats.FileSize = info.Size()
		}
	}

	return stats, nil
}

// getDefaultState 获取默认状态
func (sm *DefaultStateManager) getDefaultState() *ActiveState {
	return &ActiveState{
		CurrentSDKs: make(map[string]string),
		AddPaths:    []string{},
		AddEnvs:     make(map[string]string),
		UpdatedAt:   time.Now(),
	}
}

// StateStats 状态统计信息
type StateStats struct {
	ActiveSDKCount int       `json:"active_sdk_count"` // 激活的SDK数量
	PathCount      int       `json:"path_count"`       // 路径数量
	EnvCount       int       `json:"env_count"`        // 环境变量数量
	LastUpdated    time.Time `json:"last_updated"`     // 最后更新时间
	FileSize       int64     `json:"file_size"`        // 文件大小
}

// IsEmpty 检查状态是否为空
func (state *ActiveState) IsEmpty() bool {
	return len(state.CurrentSDKs) == 0 &&
		   len(state.AddPaths) == 0 &&
		   len(state.AddEnvs) == 0
}

// Clone 克隆状态
func (state *ActiveState) Clone() *ActiveState {
	clone := &ActiveState{
		CurrentSDKs: make(map[string]string),
		AddPaths:    make([]string, len(state.AddPaths)),
		AddEnvs:     make(map[string]string),
		UpdatedAt:   state.UpdatedAt,
	}

	for k, v := range state.CurrentSDKs {
		clone.CurrentSDKs[k] = v
	}

	copy(clone.AddPaths, state.AddPaths)

	for k, v := range state.AddEnvs {
		clone.AddEnvs[k] = v
	}

	return clone
}

// Merge 合并另一个状态
func (state *ActiveState) Merge(other *ActiveState) {
	if other == nil {
		return
	}

	// 合并SDK
	for k, v := range other.CurrentSDKs {
		state.CurrentSDKs[k] = v
	}

	// 合并路径（去重）
	pathSet := make(map[string]bool)
	for _, path := range state.AddPaths {
		pathSet[path] = true
	}

	for _, path := range other.AddPaths {
		if !pathSet[path] {
			state.AddPaths = append(state.AddPaths, path)
			pathSet[path] = true
		}
	}

	// 合并环境变量
	for k, v := range other.AddEnvs {
		state.AddEnvs[k] = v
	}

	// 更新时间戳
	if other.UpdatedAt.After(state.UpdatedAt) {
		state.UpdatedAt = other.UpdatedAt
	}
}
