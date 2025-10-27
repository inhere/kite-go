package shell

import (
	"fmt"
	"strings"
)

// ShellType shell类型枚举 bash, zsh, pwsh, cmd
type ShellType string

// ProfilePath get shell profile path
func (st ShellType) ProfilePath() string {
	switch st {
	case Bash:
		return "~/.bashrc"
	case Zsh:
		return "~/.zshrc"
	case Pwsh:
		// echo $PROFILE.CurrentUserAllHosts
		// v5: path-to-users\Documents\WindowsPowerShell\profile.ps1
		// v7: path-to-users\Documents\PowerShell\profile.ps1

		return "~/.pwsh/profile.ps1"
	case Cmd:
		// clink info:
		// C:\Users\{username}\AppData\Local\clink\ 创建 profile.lua
		return "~/AppData/Local/clink/profile.lua"
	default:
		panic("unsupported shell type: " + string(st))
	}
}

const (
	Bash ShellType = "bash"
	Zsh  ShellType = "zsh"
	Pwsh ShellType = "pwsh"
	Cmd  ShellType = "cmd"

	// Unknown shell type
	Unknown ShellType = "unknown"
)

// ScriptMark 输出的脚本必须添加标记，前面部分为message, 后面部分为脚本
const ScriptMark = "--Expression--"

var (
	// AllShellTypes 所有支持的shell类型
	AllShellTypes = []ShellType{Bash, Zsh, Pwsh, Cmd}
)

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
		return "", fmt.Errorf("unsupported shell type: %s", shellType)
	}
}
