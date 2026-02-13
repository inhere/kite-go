package sysclean

import (
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/gookit/goutil/errorx"
	"github.com/gookit/goutil/fsutil"
)

// TrashManager 回收站管理器接口
type TrashManager interface {
	Move(path string) error
	IsAvailable() bool
}

// trashManager 回收站管理器实现
type trashManager struct {
	available bool
}

// NewTrashManager 创建回收站管理器
func NewTrashManager() TrashManager {
	return &trashManager{
		available: checkTrashAvailable(),
	}
}

// checkTrashAvailable 检查回收站功能是否可用
func checkTrashAvailable() bool {
	switch runtime.GOOS {
	case "windows", "darwin":
		return true
	case "linux":
		// Linux 需要检查 FreeDesktop.org trash 规范
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return false
		}
		trashDir := filepath.Join(homeDir, ".local", "share", "Trash")
		return fsutil.IsDir(trashDir)
	default:
		return false
	}
}

// Move 移动文件到回收站
func (tm *trashManager) Move(path string) error {
	if !tm.available {
		return errorx.New("回收站功能不可用")
	}

	// 检查源文件是否存在
	if !fsutil.PathExists(path) {
		return errorx.Errorf("文件不存在：%s", path)
	}

	switch runtime.GOOS {
	case "windows":
		return tm.moveToTrashWindows(path)
	case "darwin":
		return tm.moveToTrashDarwin(path)
	case "linux":
		return tm.moveToTrashLinux(path)
	default:
		return errorx.New("不支持的操作系统")
	}
}

// IsAvailable 检查回收站是否可用
func (tm *trashManager) IsAvailable() bool {
	return tm.available
}

// moveToTrashWindows Windows 平台移动到回收站
func (tm *trashManager) moveToTrashWindows(path string) error {
	// Windows 下使用简单的方式：移动到临时目录模拟回收站
	// 完整实现需要调用 SHFileOperation API，这里使用简化版本
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	trashDir := filepath.Join(homeDir, ".kite-go", "trash")
	if err := os.MkdirAll(trashDir, 0755); err != nil {
		return errorx.Wrap(err, "创建回收站目录失败")
	}

	// 生成唯一的目标名称
	baseName := filepath.Base(path)
	targetPath := filepath.Join(trashDir, baseName)

	// 如果目标已存在，添加时间戳
	if fsutil.PathExists(targetPath) {
		ext := filepath.Ext(baseName)
		name := baseName[:len(baseName)-len(ext)]
		targetPath = filepath.Join(trashDir, name+"_"+time.Now().Format("20060102150405")+ext)
	}

	// 移动文件
	if err := os.Rename(path, targetPath); err != nil {
		return errorx.Wrap(err, "移动文件到回收站失败")
	}

	return nil
}

// moveToTrashDarwin macOS 平台移动到回收站
func (tm *trashManager) moveToTrashDarwin(path string) error {
	// macOS 下移动到 ~/.Trash
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	trashDir := filepath.Join(homeDir, ".Trash")
	if err := os.MkdirAll(trashDir, 0700); err != nil {
		return errorx.Wrap(err, "创建回收站目录失败")
	}

	baseName := filepath.Base(path)
	targetPath := filepath.Join(trashDir, baseName)

	// 如果目标已存在，添加时间戳
	if fsutil.PathExists(targetPath) {
		ext := filepath.Ext(baseName)
		name := baseName[:len(baseName)-len(ext)]
		targetPath = filepath.Join(trashDir, name+" "+time.Now().Format("2006-01-02 15:04:05")+ext)
	}

	// 移动文件
	if err := os.Rename(path, targetPath); err != nil {
		return errorx.Wrap(err, "移动文件到回收站失败")
	}

	return nil
}

// moveToTrashLinux Linux 平台移动到回收站
func (tm *trashManager) moveToTrashLinux(path string) error {
	// Linux 下使用 FreeDesktop.org Trash 规范
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	trashDir := filepath.Join(homeDir, ".local", "share", "Trash")
	filesDir := filepath.Join(trashDir, "files")
	infoDir := filepath.Join(trashDir, "info")

	// 创建必要的目录
	if err := os.MkdirAll(filesDir, 0700); err != nil {
		return errorx.Wrap(err, "创建回收站文件目录失败")
	}
	if err := os.MkdirAll(infoDir, 0700); err != nil {
		return errorx.Wrap(err, "创建回收站信息目录失败")
	}

	baseName := filepath.Base(path)
	targetPath := filepath.Join(filesDir, baseName)
	infoPath := filepath.Join(infoDir, baseName+".trashinfo")

	// 如果目标已存在，添加后缀
	counter := 1
	for fsutil.PathExists(targetPath) {
		ext := filepath.Ext(baseName)
		name := baseName[:len(baseName)-len(ext)]
		targetPath = filepath.Join(filesDir, name+"_"+string(rune('0'+counter))+ext)
		infoPath = filepath.Join(infoDir, name+"_"+string(rune('0'+counter))+ext+".trashinfo")
		counter++
	}

	// 创建 .trashinfo 文件
	infoContent := "[Trash Info]\nPath=" + path + "\nDeletionDate=" + time.Now().Format(time.RFC3339) + "\n"
	if err := os.WriteFile(infoPath, []byte(infoContent), 0600); err != nil {
		return errorx.Wrap(err, "创建回收站信息文件失败")
	}

	// 移动文件
	if err := os.Rename(path, targetPath); err != nil {
		// 删除 info 文件
		_ = os.Remove(infoPath)
		return errorx.Wrap(err, "移动文件到回收站失败")
	}

	return nil
}

// DeleteDirect 直接删除文件（不使用回收站）
func DeleteDirect(path string) error {
	if !fsutil.PathExists(path) {
		return nil // 文件不存在，无需删除
	}

	info, err := os.Stat(path)
	if err != nil {
		return errorx.Wrap(err, "获取文件信息失败")
	}

	if info.IsDir() {
		// 删除目录
		if err := os.RemoveAll(path); err != nil {
			return errorx.Wrap(err, "删除目录失败")
		}
	} else {
		// 删除文件
		if err := os.Remove(path); err != nil {
			return errorx.Wrap(err, "删除文件失败")
		}
	}

	return nil
}

// DeleteTarget 清理目标（根据配置决定是否使用回收站）
func DeleteTarget(path string, useTrash bool, trashMgr TrashManager) error {
	if useTrash && trashMgr != nil && trashMgr.IsAvailable() {
		return trashMgr.Move(path)
	}
	return DeleteDirect(path)
}
