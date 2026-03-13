//go:build !windows

package sshc

import (
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

// unixWinChangeHandler handles terminal window resize events on Unix systems.
// It listens for SIGWINCH signals and notifies the SSH session to resize the PTY.
type unixWinChangeHandler struct {
	stopChan chan struct{}
	sigChan  chan os.Signal
}

// doSetupWinChangeHandler creates a handler that listens for SIGWINCH signals.
// When the terminal window is resized, it sends the new dimensions to the SSH session.
func doSetupWinChangeHandler(session *ssh.Session, fd int) WinChangeHandler {
	h := &unixWinChangeHandler{
		stopChan: make(chan struct{}),
		sigChan:  make(chan os.Signal, 1),
	}

	signal.Notify(h.sigChan, syscall.SIGWINCH)

	go func() {
		for {
			select {
			case <-h.stopChan:
				return
			case <-h.sigChan:
				w, winH, _ := term.GetSize(fd)
				session.WindowChange(winH, w)
			}
		}
	}()

	return h
}

// Stop stops listening for SIGWINCH signals and releases resources.
func (h *unixWinChangeHandler) Stop() {
	close(h.stopChan)
	signal.Stop(h.sigChan)
}
