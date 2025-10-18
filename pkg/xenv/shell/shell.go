package shell

// ShellType shell类型枚举
type ShellType string

const (
	Bash ShellType = "bash"
	Zsh  ShellType = "zsh"
	Pwsh ShellType = "pwsh"
	Cmd  ShellType = "cmd"
)
