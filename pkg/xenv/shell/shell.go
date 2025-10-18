package shell

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
