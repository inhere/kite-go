package shell

// This file contains shell integration utilities

// IsValidShellType checks if a shell type is valid
func IsValidShellType(shellType string) bool {
	switch shellType {
	case string(Bash), string(Zsh), string(Pwsh):
		return true
	default:
		return false
	}
}
