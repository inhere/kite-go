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
	Desc: "Generate shell integration script contents",
	Help: `
<cyan>Config for Bash:</>
  // write to .bashrc OR .bash_profile
  eval "$(kite xenv shell --type bash)"

<cyan>Config for Zsh:</>
  // write to .zshrc OR .zsh_profile
  eval "$(kite xenv shell --type zsh)"

<cyan>Config for Pwsh:</>
  # write expr to profile. (find by: echo $Profile)
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
		shellName := shellCmdOpts.Type.String()
		if shellName == "" {
			if !shellCmdOpts.Reload || xenvHookShell == "" {
				return errorx.Err("please specify the shell type (bash, zsh, or pwsh)")
			}
			shellName = xenvHookShell
			c.Infoln("shell type using the XENV_HOOK_SHELL environment variable:", shellName)
		}

		shellType, err := shell.TypeFromString(shellName)
		if err != nil {
			return err
		}

		if err := xenv.Init(); err != nil {
			return err
		}

		generator := shell.NewScriptGenerator(shellType, config.Config())
		hookScript, err := generator.GenHookScripts()
		if err == nil {
			fmt.Print(hookScript)
		}
		return err
	},
}

// HookInitCmd the xenv hook init command
//  - 将会在 ~/.bashrc, ~/.zshrc, ~/.pwshrc 中执行注入hook脚本时，同时会调用当前命令，可以返回脚本内容自动执行
var HookInitCmd = &gcli.Command{
	Hidden: true, // This is an internal command
	Name:   "hook-init",
	Desc:   "Initialize the xenv hook script",
	Func: func(c *gcli.Command, args []string) error {
		return nil // TODO
	},
}
