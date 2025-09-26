package envmgr

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/gookit/goutil/fsutil"
)

// DefaultSDKManager 默认SDK管理器实现
type DefaultSDKManager struct {
	configManager ConfigManager
	httpClient    *http.Client
}

// NewSDKManager 创建SDK管理器
func NewSDKManager(configManager ConfigManager) *DefaultSDKManager {
	return &DefaultSDKManager{
		configManager: configManager,
		httpClient: &http.Client{
			Timeout: 30 * time.Minute, // 下载超时时间
		},
	}
}

// DownloadSDK 下载SDK
func (sm *DefaultSDKManager) DownloadSDK(sdk, version string) error {
	config, err := sm.configManager.GetSDKConfig(sdk)
	if err != nil {
		return fmt.Errorf("failed to get SDK config for %s: %w", sdk, err)
	}

	if config.InstallURL == "" {
		return fmt.Errorf("no install URL configured for SDK %s", sdk)
	}

	// 构建下载URL
	downloadURL := sm.buildDownloadURL(config.InstallURL, version)

	// 构建安装路径
	installPath := sm.buildInstallPath(config.InstallDir, version)

	// 检查是否已经存在
	if fsutil.IsDir(installPath) {
		return fmt.Errorf("SDK %s:%s already installed at %s", sdk, version, installPath)
	}

	// 创建临时下载目录
	tempDir := filepath.Join(os.TempDir(), fmt.Sprintf("kite-sdk-%s-%s", sdk, version))
	if err := fsutil.MkdirQuick(tempDir); err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// 下载文件
	tempFile := filepath.Join(tempDir, fmt.Sprintf("%s-%s.tar.gz", sdk, version))
	if err := sm.downloadFile(downloadURL, tempFile); err != nil {
		return fmt.Errorf("failed to download SDK: %w", err)
	}

	// 解压到临时目录
	extractDir := filepath.Join(tempDir, "extract")
	if err := sm.extractTarGz(tempFile, extractDir); err != nil {
		return fmt.Errorf("failed to extract SDK: %w", err)
	}

	// 移动到最终安装位置
	if err := sm.moveToInstallPath(extractDir, installPath); err != nil {
		return fmt.Errorf("failed to install SDK: %w", err)
	}

	return nil
}

// InstallSDK 安装SDK（下载并安装）
func (sm *DefaultSDKManager) InstallSDK(sdk, version string) error {
	return sm.DownloadSDK(sdk, version)
}

// UninstallSDK 卸载SDK
func (sm *DefaultSDKManager) UninstallSDK(sdk, version string) error {
	installPath := sm.GetSDKPath(sdk, version)

	if !fsutil.IsDir(installPath) {
		return fmt.Errorf("SDK %s:%s is not installed", sdk, version)
	}

	if err := os.RemoveAll(installPath); err != nil {
		return fmt.Errorf("failed to uninstall SDK %s:%s: %w", sdk, version, err)
	}

	return nil
}

// GetSDKPath 获取SDK路径
func (sm *DefaultSDKManager) GetSDKPath(sdk, version string) string {
	config, err := sm.configManager.GetSDKConfig(sdk)
	if err != nil {
		return ""
	}

	return sm.buildInstallPath(config.InstallDir, version)
}

// IsInstalled 检查SDK是否已安装
func (sm *DefaultSDKManager) IsInstalled(sdk, version string) bool {
	installPath := sm.GetSDKPath(sdk, version)
	return fsutil.IsDir(installPath)
}

// ListVersions 列出SDK的可用版本（已安装的版本）
func (sm *DefaultSDKManager) ListVersions(sdk string) ([]string, error) {
	config, err := sm.configManager.GetSDKConfig(sdk)
	if err != nil {
		return nil, err
	}

	// 获取SDK基础目录
	baseDir := filepath.Dir(sm.buildInstallPath(config.InstallDir, ""))

	if !fsutil.IsDir(baseDir) {
		return []string{}, nil
	}

	entries, err := os.ReadDir(baseDir)
	if err != nil {
		return nil, fmt.Errorf("failed to list SDK directory: %w", err)
	}

	var versions []string
	for _, entry := range entries {
		if entry.IsDir() {
			// 从目录名中提取版本号
			dirName := entry.Name()
			if strings.HasPrefix(dirName, sdk) {
				version := strings.TrimPrefix(dirName, sdk)
				if version != "" {
					versions = append(versions, version)
				}
			}
		}
	}

	return versions, nil
}

// GetSDKBinPath 获取SDK的可执行文件路径
func (sm *DefaultSDKManager) GetSDKBinPath(sdk, version string) string {
	sdkPath := sm.GetSDKPath(sdk, version)
	if sdkPath == "" {
		return ""
	}

	// 根据SDK类型确定可执行文件路径
	switch sdk {
	case "go":
		return filepath.Join(sdkPath, "bin")
	case "node":
		return filepath.Join(sdkPath, "bin")
	case "java":
		return filepath.Join(sdkPath, "bin")
	case "flutter":
		return filepath.Join(sdkPath, "bin")
	default:
		return filepath.Join(sdkPath, "bin")
	}
}

// buildDownloadURL 构建下载URL
func (sm *DefaultSDKManager) buildDownloadURL(urlTemplate, version string) string {
	url := urlTemplate

	// 替换版本号
	url = strings.ReplaceAll(url, "{version}", version)

	// 替换操作系统
	osName := runtime.GOOS
	switch osName {
	case "darwin":
		url = strings.ReplaceAll(url, "{os}", "darwin")
	case "linux":
		url = strings.ReplaceAll(url, "{os}", "linux")
	case "windows":
		url = strings.ReplaceAll(url, "{os}", "windows")
	}

	// 替换架构
	arch := runtime.GOARCH
	switch arch {
	case "amd64":
		url = strings.ReplaceAll(url, "{arch}", "amd64")
	case "arm64":
		url = strings.ReplaceAll(url, "{arch}", "arm64")
	case "386":
		url = strings.ReplaceAll(url, "{arch}", "386")
	}

	return url
}

// buildInstallPath 构建安装路径
func (sm *DefaultSDKManager) buildInstallPath(pathTemplate, version string) string {
	path := pathTemplate

	// 替换版本号
	path = strings.ReplaceAll(path, "{version}", version)

	return path
}

// downloadFile 下载文件
func (sm *DefaultSDKManager) downloadFile(url, filepath string) error {
	resp, err := sm.httpClient.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download from %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status %d: %s", resp.StatusCode, resp.Status)
	}

	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filepath, err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write file %s: %w", filepath, err)
	}

	return nil
}

// extractTarGz 解压tar.gz文件
func (sm *DefaultSDKManager) extractTarGz(tarPath, destDir string) error {
	if err := fsutil.MkdirQuick(destDir); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	file, err := os.Open(tarPath)
	if err != nil {
		return fmt.Errorf("failed to open tar file: %w", err)
	}
	defer file.Close()

	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read tar header: %w", err)
		}

		target := filepath.Join(destDir, header.Name)

		// 安全检查：防止路径遍历攻击
		if !strings.HasPrefix(target, destDir) {
			return fmt.Errorf("invalid file path in archive: %s", header.Name)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := fsutil.MkdirQuick(target); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", target, err)
			}
		case tar.TypeReg:
			if err := sm.extractFile(tarReader, target, header.FileInfo().Mode()); err != nil {
				return fmt.Errorf("failed to extract file %s: %w", target, err)
			}
		}
	}

	return nil
}

// extractFile 提取单个文件
func (sm *DefaultSDKManager) extractFile(reader io.Reader, path string, mode os.FileMode) error {
	// 确保目录存在
	if err := fsutil.MkdirQuick(filepath.Dir(path)); err != nil {
		return fmt.Errorf("failed to create parent directory: %w", err)
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, reader)
	if err != nil {
		return fmt.Errorf("failed to write file content: %w", err)
	}

	return nil
}

// moveToInstallPath 移动到安装路径
func (sm *DefaultSDKManager) moveToInstallPath(extractDir, installPath string) error {
	// 确保安装目录的父目录存在
	if err := fsutil.MkdirQuick(filepath.Dir(installPath)); err != nil {
		return fmt.Errorf("failed to create install parent directory: %w", err)
	}

	// 查找提取目录中的内容
	entries, err := os.ReadDir(extractDir)
	if err != nil {
		return fmt.Errorf("failed to read extract directory: %w", err)
	}

	if len(entries) == 1 && entries[0].IsDir() {
		// 如果只有一个目录，直接重命名
		sourcePath := filepath.Join(extractDir, entries[0].Name())
		if err := os.Rename(sourcePath, installPath); err != nil {
			return fmt.Errorf("failed to move extracted directory: %w", err)
		}
	} else {
		// 如果有多个文件/目录，创建安装目录并移动所有内容
		if err := fsutil.MkdirQuick(installPath); err != nil {
			return fmt.Errorf("failed to create install directory: %w", err)
		}

		for _, entry := range entries {
			sourcePath := filepath.Join(extractDir, entry.Name())
			destPath := filepath.Join(installPath, entry.Name())

			if err := os.Rename(sourcePath, destPath); err != nil {
				return fmt.Errorf("failed to move %s: %w", entry.Name(), err)
			}
		}
	}

	return nil
}

// GetSDKEnvVars 获取SDK需要设置的环境变量
func (sm *DefaultSDKManager) GetSDKEnvVars(sdk, version string) (map[string]string, error) {
	config, err := sm.configManager.GetSDKConfig(sdk)
	if err != nil {
		return nil, err
	}

	envVars := make(map[string]string)

	// 添加配置中定义的环境变量
	for k, v := range config.ActiveEnv {
		// 替换路径变量
		value := strings.ReplaceAll(v, "{sdk_path}", sm.GetSDKPath(sdk, version))
		value = strings.ReplaceAll(value, "{version}", version)
		envVars[k] = value
	}

	return envVars, nil
}

// ValidateSDK 验证SDK安装
func (sm *DefaultSDKManager) ValidateSDK(sdk, version string) error {
	if !sm.IsInstalled(sdk, version) {
		return fmt.Errorf("SDK %s:%s is not installed", sdk, version)
	}

	// installPath := sm.GetSDKPath(sdk, version)
	binPath := sm.GetSDKBinPath(sdk, version)

	// 检查bin目录是否存在
	if !fsutil.IsDir(binPath) {
		return fmt.Errorf("SDK %s:%s bin directory not found: %s", sdk, version, binPath)
	}

	// 检查主要可执行文件是否存在
	var mainExecutable string
	switch sdk {
	case "go":
		mainExecutable = "go"
	case "node":
		mainExecutable = "node"
	case "java":
		mainExecutable = "java"
	case "flutter":
		mainExecutable = "flutter"
	}

	if mainExecutable != "" {
		execPath := filepath.Join(binPath, mainExecutable)
		if runtime.GOOS == "windows" {
			execPath += ".exe"
		}

		if !fsutil.IsFile(execPath) {
			return fmt.Errorf("SDK %s:%s main executable not found: %s", sdk, version, execPath)
		}
	}

	return nil
}
