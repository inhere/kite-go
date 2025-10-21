package shell

import (
	"fmt"
	"strings"
)

// ShellType shell类型枚举
type ShellType string

var (
	// AllShellTypes 所有支持的shell类型
	AllShellTypes = []ShellType{Bash, Zsh, Pwsh, Cmd}
)

const (
	Bash ShellType = "bash"
	Zsh  ShellType = "zsh"
	Pwsh ShellType = "pwsh"
	Cmd  ShellType = "cmd"
)

// ScriptMark 输出的脚本必须添加标记，前面部分为message, 后面部分为脚本
const ScriptMark = "--Expression--"

// TypeFromString returns the shell type from a string
func TypeFromString(shellType string) (ShellType, error) {
	shellType = strings.ToLower(shellType)
	switch shellType {
	case "bash":
		return Bash, nil
	case "zsh":
		return Zsh, nil
	case "pwsh", "powershell":
		return Pwsh, nil
	case "cmd", "clink":
		return Cmd, nil
	default:
		return "", fmt.Errorf("unsupported shell type: %s (should: bash, zsh, pwsh or cmd)", shellType)
	}
}
