package shell

// This file contains shell integration utilities

// ShellType represents different shell types
type ShellType string

const (
	Bash  ShellType = "bash"
	Zsh   ShellType = "zsh"
	Pwsh  ShellType = "pwsh"
)

// IsValidShellType checks if a shell type is valid
func IsValidShellType(shellType string) bool {
	switch shellType {
	case string(Bash), string(Zsh), string(Pwsh):
		return true
	default:
		return false
	}
}