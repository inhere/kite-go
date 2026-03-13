package sshcmd

import (
	"fmt"
	"time"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/envutil"
	"github.com/inhere/kite-go/internal/pkg/sshc"
)

type sshClientOptions struct {
	username   string
	password   string
	privateKey string
	sshPort    int
	timeout    int
	keepalive  int
	strictHost bool
}

func (o *sshClientOptions) getUserName() string {
	if o.username != "" {
		return o.username
	}
	return envutil.Getenv("SSH_USER", "root")
}

func (o *sshClientOptions) getPassword() string {
	if o.password != "" {
		return o.password
	}
	return envutil.Getenv("SSH_PWD", "")
}

func NewSshClientCmd() *gcli.Command {
	opts := &sshClientOptions{
		username:   "",
		password:   "",
		sshPort:    22,
		timeout:    30,
		keepalive:  30,
		strictHost: false,
	}

	return &gcli.Command{
		Name:    "sshc",
		Desc:    "Connect to a remote SSH server with interactive session",
		Aliases: []string{"ssh-client", "sc"},
		Config: func(c *gcli.Command) {
			c.StrOpt(&opts.username, "user", "u", "", "SSH login username (or set SSH_USER env var)")
			c.StrOpt(&opts.password, "pwd", "p", "", "SSH login password (or set SSH_PWD env var)")
			c.StrOpt(&opts.privateKey, "key", "i", "", "Private key file path for SSH authentication")
			c.IntOpt(&opts.sshPort, "port", "", 22, "SSH server port")
			c.IntOpt(&opts.timeout, "timeout", "t", 30, "SSH connection timeout in seconds")
			c.IntOpt(&opts.keepalive, "keepalive", "k", 30, "Keepalive interval in seconds (0 to disable)")
			c.BoolOpt(&opts.strictHost, "strict", "", false, "Enable strict host key checking")

			c.AddArg("host", "SSH server hostname or IP address", true)
		},
		Func: func(c *gcli.Command, args []string) error {
			return runSshClient(c, opts, args)
		},
	}
}

func runSshClient(c *gcli.Command, opts *sshClientOptions, args []string) error {
	host := c.Arg("host").String()
	if host == "" {
		return fmt.Errorf("host is required")
	}

	username := opts.getUserName()
	password := opts.getPassword()

	if password == "" && opts.privateKey == "" {
		return fmt.Errorf("either password (--pwd) or private key (--key) is required")
	}

	addr := fmt.Sprintf("%s:%d", host, opts.sshPort)
	c.Infof("Connecting to %s as user '%s'...\n", addr, username)

	client := sshc.NewClient(username, password)
	client.Port = opts.sshPort
	client.PrivateKey = opts.privateKey
	client.Timeout = time.Duration(opts.timeout) * time.Second
	client.Keepalive = time.Duration(opts.keepalive) * time.Second
	client.StrictHost = opts.strictHost

	if err := client.Connect(host); err != nil {
		return err
	}
	defer client.Close()

	c.Successf("Connected to %s\n", addr)

	return client.StartInteractive()
}
