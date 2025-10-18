package subcmd

import (
	"fmt"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/gflag"
	"github.com/gookit/goutil/cflag"
	"github.com/gookit/goutil/envutil"
	"github.com/gookit/goutil/errorx"
	"github.com/inhere/kite-go/pkg/xenv/config"
	"github.com/inhere/kite-go/pkg/xenv/shell"
)

var shellCmdOpts = struct {
	Type gflag.EnumString
	Reload  bool
}{
	Type: cflag.NewEnumString("bash", "zsh", "pwsh"),
}

// ShellCmd the xenv shell command
var ShellCmd = &gcli.Command{
	Name: "shell",
	Desc: "Generate shell integration script",
	Config: func(c *gcli.Command) {
		c.BoolOpt(&shellCmdOpts.Reload, "reload", "r", false, "Reload the xenv shell script codes")
		c.VarOpt(&shellCmdOpts.Type, "type", "t", "Shell type (bash, zsh, pwsh)")
	},
	Func: func(c *gcli.Command, args []string) error {
		// XENV_HOOK_SHELL for reload shell hook script
		xenvHookShell := envutil.Getenv("XENV_HOOK_SHELL")
		shellType := shellCmdOpts.Type.String()
		if shellType == "" {
			if !shellCmdOpts.Reload || xenvHookShell == "" {
				return errorx.Err("please specify the shell type (bash, zsh, or pwsh)")
			}
			shellType = xenvHookShell
			c.Infoln("shell type using the XENV_HOOK_SHELL environment variable:", shellType)
		}

		if err := config.Mgr.Init(); err != nil {
			return err
		}

		generator := shell.NewScriptGenerator(config.Config())
		hookScript, err := generator.GenerateScripts(shellType)
		if err == nil {
			fmt.Print(hookScript)
		}
		return err
	},
}
