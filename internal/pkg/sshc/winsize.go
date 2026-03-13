package sshc

import (
	"golang.org/x/crypto/ssh"
)

// WinChangeHandler is an interface for handling terminal window resize events.
// Implementations are platform-specific:
//   - On Unix systems: handles SIGWINCH signals to resize the PTY
//   - On Windows: no-op since SIGWINCH is not available
type WinChangeHandler interface {
	// Stop stops listening for window change events and releases resources.
	Stop()
}

// setupWinChangeHandler creates a platform-specific window change handler.
// On Unix systems, it listens for SIGWINCH signals and notifies the SSH session
// when the terminal size changes. On Windows, it returns a no-op handler.
//
// This is an internal function used by StartInteractiveSession.
func setupWinChangeHandler(session *ssh.Session, fd int) WinChangeHandler {
	return doSetupWinChangeHandler(session, fd)
}
