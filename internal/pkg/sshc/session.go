package sshc

import (
	"fmt"
	"os"

	"github.com/gookit/goutil/x/ccolor"
	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

// StartInteractiveSession starts an interactive terminal session using the provided SSH client.
// It performs the following:
//   - Creates a new SSH session
//   - Sets the local terminal to raw mode
//   - Requests a PTY (pseudo-terminal) with the current terminal size
//   - Redirects stdin/stdout/stderr to the session
//   - Handles terminal resize events (on Unix systems)
//   - Waits for the session to end
//
// The function blocks until the session ends (user types 'exit', presses Ctrl+D,
// or the connection is closed). After the session ends, it restores the terminal
// to its original state.
//
// Note: On Windows, terminal resize events are not supported.
//
// Example:
//
//	client, err := ssh.Dial("tcp", "host:22", config)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer client.Close()
//
//	if err := sshc.StartInteractiveSession(client); err != nil {
//	    log.Fatal(err)
//	}
func StartInteractiveSession(client *ssh.Client) error {
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()

	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return fmt.Errorf("failed to set terminal to raw mode: %w", err)
	}
	defer term.Restore(fd, oldState)

	w, h, err := term.GetSize(fd)
	if err != nil {
		w, h = 80, 24
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	if err := session.RequestPty("xterm-256color", h, w, modes); err != nil {
		return fmt.Errorf("failed to request PTY: %w", err)
	}

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	if err := session.Shell(); err != nil {
		return fmt.Errorf("failed to start shell: %w", err)
	}

	winChangeHandler := setupWinChangeHandler(session, fd)
	defer winChangeHandler.Stop()

	ccolor.Infoln("Session started. Press Ctrl+D or type 'exit' to disconnect.")

	if err := session.Wait(); err != nil {
		if exitErr, ok := err.(*ssh.ExitError); ok {
			ccolor.Fprintf(os.Stderr, "<yellow>Session ended with exit code:</> %d\n", exitErr.ExitStatus())
		}
	}

	ccolor.Infoln("Session closed.")
	return nil
}
