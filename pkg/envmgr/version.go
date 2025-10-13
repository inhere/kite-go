package envmgr

import (
	"errors"
	"strings"
)

var (
	ErrInvalidVersionSpec = errors.New("invalid version specification")
	ErrEmptyVersionSpec   = errors.New("empty version specification")
)

// ParseVersionSpec 解析版本规格 "sdk" or "sdk:version" or "sdk@version"
func ParseVersionSpec(spec string) (*VersionSpec, error) {
	if spec == "" {
		return nil, ErrEmptyVersionSpec
	}

	sep := ":"
	if strings.Contains(spec, "@") {
		sep = "@"
	}

	parts := strings.SplitN(spec, sep, 2)
	if len(parts) != 2 {
		parts = append(parts, "latest")
	}

	sdk := strings.TrimSpace(parts[0])
	version := strings.TrimSpace(parts[1])
	if sdk == "" || version == "" {
		return nil, ErrInvalidVersionSpec
	}

	return &VersionSpec{
		SDK:     sdk,
		Version: version,
	}, nil
}

// ParseMultipleVersionSpecs 解析多个版本规格
func ParseMultipleVersionSpecs(specs []string) ([]*VersionSpec, error) {
	if len(specs) == 0 {
		return nil, nil
	}

	var result []*VersionSpec
	for _, spec := range specs {
		parsed, err := ParseVersionSpec(spec)
		if err != nil {
			return nil, err
		}
		result = append(result, parsed)
	}

	return result, nil
}

// String 返回版本规格的字符串表示
func (vs *VersionSpec) String() string {
	return vs.SDK + ":" + vs.Version
}

// IsValidSDKName 检查SDK名称是否有效
func IsValidSDKName(name string) bool {
	if name == "" {
		return false
	}

	// SDK名称只能包含字母、数字、下划线和连字符
	for _, r := range name {
		if !((r >= 'a' && r <= 'z') ||
			(r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') ||
			 r == '_' || r == '-') {
			return false
		}
	}

	return true
}

// IsValidVersion 检查版本号是否有效
func IsValidVersion(version string) bool {
	if version == "" {
		return false
	}

	// 支持的版本格式：
	// 1. 语义版本: 1.2.3, 1.2.3-alpha, 1.2.3+build
	// 2. 主版本: 18, 16
	// 3. 别名: lts, latest, stable
	// 4. 自动检测: auto

	// 简单验证：不能包含空格和特殊字符
	invalidChars := []string{" ", "\t", "\n", "\r"}
	for _, char := range invalidChars {
		if strings.Contains(version, char) {
			return false
		}
	}

	return true
}

// NormalizeVersion 标准化版本号
func NormalizeVersion(version string) string {
	version = strings.TrimSpace(version)

	// 处理别名
	switch strings.ToLower(version) {
	case "lts":
		return "lts"
	case "latest":
		return "latest"
	case "stable":
		return "stable"
	case "auto":
		return "auto"
	default:
		return version
	}
}

// CompareVersions 比较两个版本号
// 返回值: -1 表示 v1 < v2, 0 表示 v1 == v2, 1 表示 v1 > v2
func CompareVersions(v1, v2 string) int {
	// 简单的字符串比较，实际应该实现语义版本比较
	if v1 == v2 {
		return 0
	}
	if v1 < v2 {
		return -1
	}
	return 1
}
