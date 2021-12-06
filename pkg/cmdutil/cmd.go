package cmdutil

import (
	"os"
	"os/exec"

	"github.com/gookit/color"
	"github.com/gookit/goutil/cliutil"
)

// Cmd struct
type Cmd struct {
	Bin  string
	Args []string
	// line = bin + args
	line string
	// hooks
	Before func() bool
}

// NewCmd new command.
func NewCmd() *Cmd {
	return &Cmd{}
}

func NewGitCmd(subCmd string, args ...string) *Cmd {
	c := &Cmd{
		Bin:  "git",
		Args: []string{subCmd},
	}

	return c.AddArgs(args...)
}

func NewCmdWithLine(line string) *Cmd {
	return &Cmd{line: line}
}

// SetBinArgs to command
func (c *Cmd) SetBinArgs(binName string, args ...string) *Cmd {
	c.Bin = binName
	c.Args = args
	return c
}

// AddArgs to command
func (c *Cmd) AddArgs(args ...string) *Cmd {
	c.Args = append(c.Args, args...)
	return c
}

// NewExecCmd create exec.Cmd from current cmd
func (c *Cmd) NewExecCmd() *exec.Cmd {
	c.parseBinArgs()

	// create exec.Cmd
	return exec.Command(c.Bin, c.Args...)
}

// MustRun cmd
func (c *Cmd) MustRun() {
	err := c.Run()
	if err != nil {
		color.Errorln(err.Error())
	}
}

// Run cmd
func (c *Cmd) Run() error {
	c.parseBinArgs()

	// create exec.Cmd
	osc := exec.Command(c.Bin, c.Args...)

	osc.Stdout = os.Stdout
	osc.Stderr = os.Stderr

	return osc.Run()
}

// GetBinArgs cmd line string
func (c *Cmd) GetBinArgs() (string, []string) {
	c.parseBinArgs()
	return c.Bin, c.Args
}

// parse cmd line string
func (c *Cmd) parseBinArgs() {
	if c.Bin != "" {
		return
	}

	if c.line != "" {
		args := cliutil.ParseLine(c.line)
		// binding to cmd
		c.Bin = args[0]
		if len(args) > 1 {
			c.Args = args[1:]
		}
	}
}

// String cmd line string
func (c *Cmd) String() string {
	if c.line == "" {
		c.line = cliutil.LineBuild(c.Bin, c.Args)
	}

	return c.line
}
