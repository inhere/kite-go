package envmgr

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestDefaultStateManager(t *testing.T) {
	// 创建临时目录用于测试
	tempDir := t.TempDir()
	stateFile := filepath.Join(tempDir, "test_state.json")
	
	sm := NewStateManager(stateFile)

	t.Run("LoadState with non-existent file", func(t *testing.T) {
		state, err := sm.LoadState()
		if err != nil {
			t.Errorf("LoadState() failed: %v", err)
		}
		
		if state == nil {
			t.Error("LoadState() returned nil state")
		}
		
		if len(state.CurrentSDKs) != 0 {
			t.Errorf("Expected empty CurrentSDKs, got %v", state.CurrentSDKs)
		}
	})

	t.Run("SaveState and LoadState", func(t *testing.T) {
		originalState := &ActiveState{
			CurrentSDKs: map[string]string{
				"go":   "1.21.5",
				"node": "18",
			},
			AddPaths: []string{"/opt/go/bin", "/opt/node/bin"},
			AddEnvs: map[string]string{
				"GOROOT": "/opt/go",
				"NODE_ENV": "development",
			},
			UpdatedAt: time.Now(),
		}
		
		err := sm.SaveState(originalState)
		if err != nil {
			t.Errorf("SaveState() failed: %v", err)
		}
		
		// 验证文件存在
		if !sm.IsStateFileExists() {
			t.Error("State file should exist after SaveState()")
		}
		
		// 加载状态
		loadedState, err := sm.LoadState()
		if err != nil {
			t.Errorf("LoadState() failed: %v", err)
		}
		
		// 验证SDK状态
		if len(loadedState.CurrentSDKs) != len(originalState.CurrentSDKs) {
			t.Errorf("Expected %d SDKs, got %d", len(originalState.CurrentSDKs), len(loadedState.CurrentSDKs))
		}
		
		for sdk, version := range originalState.CurrentSDKs {
			if loadedState.CurrentSDKs[sdk] != version {
				t.Errorf("Expected %s version %s, got %s", sdk, version, loadedState.CurrentSDKs[sdk])
			}
		}
		
		// 验证路径
		if len(loadedState.AddPaths) != len(originalState.AddPaths) {
			t.Errorf("Expected %d paths, got %d", len(originalState.AddPaths), len(loadedState.AddPaths))
		}
		
		// 验证环境变量
		if len(loadedState.AddEnvs) != len(originalState.AddEnvs) {
			t.Errorf("Expected %d env vars, got %d", len(originalState.AddEnvs), len(loadedState.AddEnvs))
		}
	})

	t.Run("UpdateSDKState", func(t *testing.T) {
		err := sm.UpdateSDKState("python", "3.11", true)
		if err != nil {
			t.Errorf("UpdateSDKState() failed: %v", err)
		}
		
		state, err := sm.LoadState()
		if err != nil {
			t.Errorf("LoadState() failed: %v", err)
		}
		
		if state.CurrentSDKs["python"] != "3.11" {
			t.Errorf("Expected python version 3.11, got %s", state.CurrentSDKs["python"])
		}
		
		// 移除SDK
		err = sm.UpdateSDKState("python", "", false)
		if err != nil {
			t.Errorf("UpdateSDKState() failed: %v", err)
		}
		
		state, err = sm.LoadState()
		if err != nil {
			t.Errorf("LoadState() failed: %v", err)
		}
		
		if _, exists := state.CurrentSDKs["python"]; exists {
			t.Error("Expected python to be removed from CurrentSDKs")
		}
	})

	t.Run("AddPath and RemovePath", func(t *testing.T) {
		testPath := "/test/path"
		
		err := sm.AddPath(testPath)
		if err != nil {
			t.Errorf("AddPath() failed: %v", err)
		}
		
		state, err := sm.LoadState()
		if err != nil {
			t.Errorf("LoadState() failed: %v", err)
		}
		
		found := false
		for _, path := range state.AddPaths {
			if path == testPath {
				found = true
				break
			}
		}
		
		if !found {
			t.Errorf("Expected path %s to be added", testPath)
		}
		
		// 移除路径
		err = sm.RemovePath(testPath)
		if err != nil {
			t.Errorf("RemovePath() failed: %v", err)
		}
		
		state, err = sm.LoadState()
		if err != nil {
			t.Errorf("LoadState() failed: %v", err)
		}
		
		for _, path := range state.AddPaths {
			if path == testPath {
				t.Errorf("Expected path %s to be removed", testPath)
			}
		}
	})

	t.Run("SetEnv and UnsetEnv", func(t *testing.T) {
		testEnvName := "TEST_VAR"
		testEnvValue := "test_value"
		
		err := sm.SetEnv(testEnvName, testEnvValue)
		if err != nil {
			t.Errorf("SetEnv() failed: %v", err)
		}
		
		state, err := sm.LoadState()
		if err != nil {
			t.Errorf("LoadState() failed: %v", err)
		}
		
		if state.AddEnvs[testEnvName] != testEnvValue {
			t.Errorf("Expected env %s=%s, got %s", testEnvName, testEnvValue, state.AddEnvs[testEnvName])
		}
		
		// 移除环境变量
		err = sm.UnsetEnv(testEnvName)
		if err != nil {
			t.Errorf("UnsetEnv() failed: %v", err)
		}
		
		state, err = sm.LoadState()
		if err != nil {
			t.Errorf("LoadState() failed: %v", err)
		}
		
		if _, exists := state.AddEnvs[testEnvName]; exists {
			t.Errorf("Expected env %s to be removed", testEnvName)
		}
	})

	t.Run("GetStateStats", func(t *testing.T) {
		stats, err := sm.GetStateStats()
		if err != nil {
			t.Errorf("GetStateStats() failed: %v", err)
		}
		
		if stats == nil {
			t.Error("GetStateStats() returned nil")
		}
		
		if stats.FileSize <= 0 {
			t.Error("Expected file size to be greater than 0")
		}
	})

	t.Run("BackupState", func(t *testing.T) {
		err := sm.BackupState()
		if err != nil {
			t.Errorf("BackupState() failed: %v", err)
		}
		
		// 检查备份文件是否存在
		backupPattern := stateFile + ".backup.*"
		matches, err := filepath.Glob(backupPattern)
		if err != nil {
			t.Errorf("Failed to check backup files: %v", err)
		}
		
		if len(matches) == 0 {
			t.Error("Expected backup file to be created")
		}
	})
}

func TestActiveState(t *testing.T) {
	t.Run("IsEmpty", func(t *testing.T) {
		emptyState := &ActiveState{
			CurrentSDKs: make(map[string]string),
			AddPaths:    []string{},
			AddEnvs:     make(map[string]string),
		}
		
		if !emptyState.IsEmpty() {
			t.Error("Expected empty state to return true for IsEmpty()")
		}
		
		nonEmptyState := &ActiveState{
			CurrentSDKs: map[string]string{"go": "1.21"},
			AddPaths:    []string{},
			AddEnvs:     make(map[string]string),
		}
		
		if nonEmptyState.IsEmpty() {
			t.Error("Expected non-empty state to return false for IsEmpty()")
		}
	})

	t.Run("Clone", func(t *testing.T) {
		original := &ActiveState{
			CurrentSDKs: map[string]string{"go": "1.21", "node": "18"},
			AddPaths:    []string{"/opt/go/bin", "/opt/node/bin"},
			AddEnvs:     map[string]string{"GOROOT": "/opt/go"},
			UpdatedAt:   time.Now(),
		}
		
		clone := original.Clone()
		
		// 修改原始状态
		original.CurrentSDKs["python"] = "3.11"
		original.AddPaths = append(original.AddPaths, "/opt/python/bin")
		original.AddEnvs["PYTHON_HOME"] = "/opt/python"
		
		// 验证克隆未受影响
		if len(clone.CurrentSDKs) != 2 {
			t.Errorf("Expected clone to have 2 SDKs, got %d", len(clone.CurrentSDKs))
		}
		
		if len(clone.AddPaths) != 2 {
			t.Errorf("Expected clone to have 2 paths, got %d", len(clone.AddPaths))
		}
		
		if len(clone.AddEnvs) != 1 {
			t.Errorf("Expected clone to have 1 env var, got %d", len(clone.AddEnvs))
		}
	})

	t.Run("Merge", func(t *testing.T) {
		state1 := &ActiveState{
			CurrentSDKs: map[string]string{"go": "1.21"},
			AddPaths:    []string{"/opt/go/bin"},
			AddEnvs:     map[string]string{"GOROOT": "/opt/go"},
			UpdatedAt:   time.Now().Add(-time.Hour),
		}
		
		state2 := &ActiveState{
			CurrentSDKs: map[string]string{"node": "18", "go": "1.22"}, // go版本会被覆盖
			AddPaths:    []string{"/opt/node/bin", "/opt/go/bin"},       // 重复路径不会添加
			AddEnvs:     map[string]string{"NODE_ENV": "production"},
			UpdatedAt:   time.Now(),
		}
		
		state1.Merge(state2)
		
		// 验证SDK合并
		if len(state1.CurrentSDKs) != 2 {
			t.Errorf("Expected 2 SDKs after merge, got %d", len(state1.CurrentSDKs))
		}
		
		if state1.CurrentSDKs["go"] != "1.22" {
			t.Errorf("Expected go version to be overwritten to 1.22, got %s", state1.CurrentSDKs["go"])
		}
		
		// 验证路径去重
		if len(state1.AddPaths) != 2 {
			t.Errorf("Expected 2 unique paths after merge, got %d", len(state1.AddPaths))
		}
		
		// 验证环境变量合并
		if len(state1.AddEnvs) != 2 {
			t.Errorf("Expected 2 env vars after merge, got %d", len(state1.AddEnvs))
		}
		
		// 验证时间戳更新
		if !state1.UpdatedAt.Equal(state2.UpdatedAt) {
			t.Error("Expected UpdatedAt to be updated to newer timestamp")
		}
	})
}