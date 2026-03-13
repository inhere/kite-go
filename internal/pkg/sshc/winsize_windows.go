//go:build windows

package sshc

import (
	"golang.org/x/crypto/ssh"
)

// windowsWinChangeHandler is a no-op implementation for Windows systems.
// Windows does not support SIGWINCH signal, so window resize events are not handled.
type windowsWinChangeHandler struct{}

// doSetupWinChangeHandler returns a no-op handler for Windows systems.
// Since Windows does not support SIGWINCH, this handler does nothing.
func doSetupWinChangeHandler(_ *ssh.Session, _ int) WinChangeHandler {
	return &windowsWinChangeHandler{}
}

// Stop is a no-op for Windows systems.
func (h *windowsWinChangeHandler) Stop() {
}
