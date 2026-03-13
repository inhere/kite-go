package sshc

import (
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
)

// Client represents an SSH client connection with configuration options.
// It provides methods for connecting to remote servers, executing commands,
// starting interactive sessions, and port forwarding.
type Client struct {
	// Username is the SSH login username
	Username string
	// Password is the SSH login password (used for password auth or as passphrase for private key)
	Password string
	// PrivateKey is the path to the private key file for key-based authentication
	PrivateKey string
	// Port is the SSH server port (default: 22)
	Port int
	// Timeout is the connection timeout duration (default: 30s)
	Timeout time.Duration
	// Keepalive is the interval for sending keepalive packets (default: 30s, 0 to disable)
	Keepalive time.Duration
	// StrictHost enables strict host key checking when set to true
	StrictHost bool

	client *ssh.Client
}

// NewClient creates a new SSH client with the given username and password.
// Default values: Port=22, Timeout=30s, Keepalive=30s, StrictHost=false.
//
// Example:
//
//	client := sshc.NewClient("root", "password")
//	client.Port = 2222
//	client.Timeout = 60 * time.Second
//	if err := client.Connect("192.168.1.100"); err != nil {
//	    log.Fatal(err)
//	}
//	defer client.Close()
func NewClient(username, password string) *Client {
	return &Client{
		Username:   username,
		Password:   password,
		Port:       22,
		Timeout:    30 * time.Second,
		Keepalive:  30 * time.Second,
		StrictHost: false,
	}
}

// Connect establishes a TCP connection to the SSH server at the given host.
// It authenticates using the configured credentials (password and/or private key).
// If Keepalive is set, it starts a background goroutine to send keepalive packets.
//
// The host parameter should be the hostname or IP address without port (use Client.Port for port).
func (c *Client) Connect(host string) error {
	addr := fmt.Sprintf("%s:%d", host, c.Port)

	config := &ssh.ClientConfig{
		User:            c.Username,
		Timeout:         c.Timeout,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	authMethods, err := c.buildAuthMethods()
	if err != nil {
		return err
	}
	config.Auth = authMethods

	if c.StrictHost {
		config.HostKeyCallback = ssh.FixedHostKey(nil)
	}

	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return fmt.Errorf("failed to connect to %s: %w", addr, err)
	}

	c.client = client

	if c.Keepalive > 0 {
		go c.startKeepalive()
	}

	return nil
}

func (c *Client) buildAuthMethods() ([]ssh.AuthMethod, error) {
	var methods []ssh.AuthMethod

	if c.PrivateKey != "" {
		keyContent, err := os.ReadFile(c.PrivateKey)
		if err != nil {
			return nil, fmt.Errorf("failed to read private key file: %w", err)
		}

		var signer ssh.Signer
		if c.Password != "" {
			signer, err = ssh.ParsePrivateKeyWithPassphrase(keyContent, []byte(c.Password))
		} else {
			signer, err = ssh.ParsePrivateKey(keyContent)
		}
		if err != nil {
			return nil, fmt.Errorf("failed to parse private key: %w", err)
		}
		methods = append(methods, ssh.PublicKeys(signer))
	}

	if c.Password != "" {
		methods = append(methods, ssh.Password(c.Password))
	}

	return methods, nil
}

func (c *Client) startKeepalive() {
	ticker := time.NewTicker(c.Keepalive)
	defer ticker.Stop()

	for range ticker.C {
		if c.client == nil {
			return
		}
		_, _, err := c.client.SendRequest("keepalive@golang.org", true, nil)
		if err != nil {
			return
		}
	}
}

// Close closes the SSH connection and releases resources.
// It is safe to call Close multiple times.
func (c *Client) Close() error {
	if c.client != nil {
		return c.client.Close()
	}
	return nil
}

// RawClient returns the underlying ssh.Client for advanced usage.
// Returns nil if not connected.
func (c *Client) RawClient() *ssh.Client {
	return c.client
}

// NewSession creates a new SSH session on the established connection.
// The caller is responsible for closing the session when done.
//
// Returns an error if the client is not connected.
func (c *Client) NewSession() (*ssh.Session, error) {
	if c.client == nil {
		return nil, fmt.Errorf("not connected")
	}
	return c.client.NewSession()
}

// Execute runs a command on the remote server and returns the combined output.
// It creates a new session, runs the command, and closes the session automatically.
//
// Example:
//
//	output, err := client.Execute("ls -la /home")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(output)
func (c *Client) Execute(cmd string) (string, error) {
	session, err := c.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	output, err := session.CombinedOutput(cmd)
	return string(output), err
}

// StartInteractive starts an interactive terminal session with the remote server.
// It sets up PTY (pseudo-terminal) with the current terminal size and handles
// window resize events on Unix systems.
//
// The function blocks until the session ends (user exits or connection closes).
//
// Returns an error if the client is not connected.
func (c *Client) StartInteractive() error {
	if c.client == nil {
		return fmt.Errorf("not connected")
	}
	return StartInteractiveSession(c.client)
}

// ForwardLocal creates a local port forward from localAddr to remoteAddr.
// It listens on localAddr and forwards connections to remoteAddr through the SSH tunnel.
//
// This is a blocking function that runs indefinitely until an error occurs.
// For typical use, run it in a goroutine.
//
// Example:
//
//	// Forward local port 8080 to remote port 80
//	go client.ForwardLocal("localhost:8080", "localhost:80")
func (c *Client) ForwardLocal(localAddr, remoteAddr string) error {
	if c.client == nil {
		return fmt.Errorf("not connected")
	}

	listener, err := net.Listen("tcp", localAddr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", localAddr, err)
	}
	defer listener.Close()

	for {
		localConn, err := listener.Accept()
		if err != nil {
			return fmt.Errorf("failed to accept connection: %w", err)
		}

		remoteConn, err := c.client.Dial("tcp", remoteAddr)
		if err != nil {
			localConn.Close()
			continue
		}

		go func() {
			defer localConn.Close()
			defer remoteConn.Close()
			io.Copy(localConn, remoteConn)
		}()

		go func() {
			defer localConn.Close()
			defer remoteConn.Close()
			io.Copy(remoteConn, localConn)
		}()
	}
}

// Dial establishes a connection to the remote address through the SSH tunnel.
// It uses the SSH client's Dial method to create the connection.
//
// The network parameter should be "tcp" or "udp".
// The addr parameter should be in the format "host:port".
//
// Returns an error if the client is not connected.
func (c *Client) Dial(network, addr string) (net.Conn, error) {
	if c.client == nil {
		return nil, fmt.Errorf("not connected")
	}
	return c.client.Dial(network, addr)
}
