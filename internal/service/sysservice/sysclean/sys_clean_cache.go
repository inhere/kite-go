package sysclean

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/gookit/goutil/errorx"
)

// CacheManager 缓存管理器
type CacheManager struct {
	filePath string
	ttl      time.Duration
}

// NewCacheManager 创建缓存管理器
func NewCacheManager(filePath string, ttl time.Duration) *CacheManager {
	if ttl <= 0 {
		ttl = 3 * time.Minute
	}
	return &CacheManager{
		filePath: filePath,
		ttl:      ttl,
	}
}

// Load 加载缓存
func (cm *CacheManager) Load() (*ScanResult, error) {
	data, err := os.ReadFile(cm.filePath)
	if err != nil {
		return nil, errorx.Wrap(err, "读取缓存文件失败")
	}

	var result ScanResult
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, errorx.Wrap(err, "解析缓存数据失败")
	}

	// 检查是否过期
	if time.Now().After(result.ExpiresAt) {
		return nil, errorx.New("缓存已过期")
	}

	return &result, nil
}

// Save 保存缓存
func (cm *CacheManager) Save(result *ScanResult) error {
	// 设置过期时间
	result.CreatedAt = time.Now()
	result.ExpiresAt = result.CreatedAt.Add(cm.ttl)

	// 确保目录存在
	dir := filepath.Dir(cm.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return errorx.Wrap(err, "创建缓存目录失败")
	}

	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return errorx.Wrap(err, "序列化缓存数据失败")
	}

	if err := os.WriteFile(cm.filePath, data, 0644); err != nil {
		return errorx.Wrap(err, "写入缓存文件失败")
	}

	return nil
}

// Clear 清除缓存
func (cm *CacheManager) Clear() error {
	if _, err := os.Stat(cm.filePath); os.IsNotExist(err) {
		return nil // 文件不存在，无需清除
	}

	if err := os.Remove(cm.filePath); err != nil {
		return errorx.Wrap(err, "删除缓存文件失败")
	}

	return nil
}

// IsExpired 检查缓存是否过期
func (cm *CacheManager) IsExpired() bool {
	info, err := os.Stat(cm.filePath)
	if err != nil {
		return true
	}

	// 根据文件的修改时间和 TTL 判断
	return time.Since(info.ModTime()) > cm.ttl
}

// RemainingTime 获取缓存剩余有效时间
func (cm *CacheManager) RemainingTime() time.Duration {
	info, err := os.Stat(cm.filePath)
	if err != nil {
		return 0
	}

	elapsed := time.Since(info.ModTime())
	remaining := cm.ttl - elapsed
	if remaining < 0 {
		return 0
	}
	return remaining
}

// Exists 检查缓存是否存在
func (cm *CacheManager) Exists() bool {
	_, err := os.Stat(cm.filePath)
	return err == nil
}

// GetFilePath 获取缓存文件路径
func (cm *CacheManager) GetFilePath() string {
	return cm.filePath
}

// GetTTL 获取缓存 TTL
func (cm *CacheManager) GetTTL() time.Duration {
	return cm.ttl
}
