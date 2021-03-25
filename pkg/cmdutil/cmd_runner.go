package cmdutil

import (
	"fmt"
	"os"

	"github.com/gookit/color"
	"github.com/gookit/gcli/v3/interact"
	"github.com/gookit/goutil/cliutil"
)

// CmdRunner struct StepRunner
type CmdRunner struct {
	wordDir string
	lastErr error
	// Dry run all commands
	DryRun bool
	// Ignore check prevision return code
	IgnoreErr bool
	// Interactive Whether to interactively ask before executing command
	Interactive bool
	// added commands
	commands  []*Cmd
}

// NewRunner create
func NewRunner(fn ...func(cr *CmdRunner)) *CmdRunner {
	cr := &CmdRunner{}

	if len(fn) > 0 {
		fn[0](cr)
	}

	return cr
}

func (r *CmdRunner) SetWordDir(wordDir string) {
	r.wordDir = wordDir
}

// Add an command
func (r *CmdRunner) Add(binName string, args ...string) *CmdRunner {
	return r.AddWithArgs(binName, args...)
}

// Addf an command
func (r *CmdRunner) Addf(cmdFmt string, args ...interface{}) *CmdRunner {
	cmd := NewCmdWithLine(fmt.Sprintf(cmdFmt, args...))

	r.commands = append(r.commands, cmd)
	return r
}

// Add an command
func (r *CmdRunner) AddLine(cmdLine string) *CmdRunner {
	r.commands = append(r.commands, NewCmdWithLine(cmdLine))
	return r
}

// AddCmd an command
func (r *CmdRunner) AddCmd(cmd *Cmd) *CmdRunner {
	r.commands = append(r.commands, cmd)
	return r
}

// NewCmd an command
func (r *CmdRunner) NewCmd(binName string, args ...string) *Cmd {
	cmd := NewCmd()
	cmd.SetBinArgs(binName, args...)
	r.AddCmd(cmd)
	return cmd
}

// NewCmd an command
func (r *CmdRunner) NewGitCmd(subCmd string, args ...string) *Cmd {
	cmd := NewCmd()
	cmd.SetBinArgs("git", subCmd)
	cmd.AddArgs(args...)

	r.AddCmd(cmd)
	return cmd
}

// AddGitCmd an command
func (r *CmdRunner) AddGitCmd(subCmd string, args ...string) *CmdRunner {
	r.NewGitCmd(subCmd, args...)
	return r
}

// AddWithArgs add command with args
func (r *CmdRunner) AddWithArgs(binName string, args ...string) *CmdRunner {
	cmdLine := cliutil.LineBuild(binName, args)
	cmd := NewCmdWithLine(cmdLine)
	cmd.Bin = binName
	cmd.Args = args

	r.commands = append(r.commands, cmd)

	return r
}

// Run all commands.
func (r *CmdRunner) Run() {
	color.Magenta.Printf("# Run All Workflows(%d steps):\n", len(r.commands))

	yesRun := true
	// c := exec.Command("test")
	// c := exec.Cmd{}
	for i, cmd := range r.commands {
		// c := exec.Command(cmd.Bin, cmd.Args...)
		c := cmd.NewExecCmd()
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr

		color.Magenta.Println("STEP", i+1)
		color.Comment.Println(">", cmd.String())

		if r.DryRun {
			color.Infoln("DRY-RUN: command execute completed")
			continue
		}

		if r.Interactive {
			yesRun = interact.Confirm("continue run", yesRun)
			if !yesRun {
				continue
			}
		}

		r.lastErr = c.Run()
		if r.lastErr != nil && r.IgnoreErr == false {
			color.Errorln("cmd exec error:", r.lastErr, ", stop run.")
			break
		}
	}
	color.Success.Println("Run Completed")
}

func (r *CmdRunner) RunNoPrint() {

}
