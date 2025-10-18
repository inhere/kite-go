package subcmd

import (
	"fmt"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/gflag"
	"github.com/gookit/goutil/cflag"
	"github.com/gookit/goutil/envutil"
	"github.com/gookit/goutil/errorx"
	"github.com/inhere/kite-go/pkg/xenv"
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
	Help: `
<mga>Config for Bash:</>
  // write to .bashrc OR .bash_profile
  eval "$(kite xenv shell --type bash)"

<mga>Config for Zsh:<mga>
  // write to .zshrc OR .zsh_profile
  eval "$(kite xenv shell --type zsh)"

<mga>Config for Pwsh:<mga>
  # write to profile. (find by: echo $Profile)
  # Method 1:
  Invoke-Expression (&kite xenv shell --type pwsh)
  # Method 2:
  kite xenv shell --type pwsh | Out-String | Invoke-Expression
`,
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

		if err := xenv.Init(); err != nil {
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
