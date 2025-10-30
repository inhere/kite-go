package subcmd

import (
	"fmt"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/gflag"
	"github.com/gookit/goutil/cflag"
	"github.com/gookit/goutil/errorx"
	"github.com/gookit/goutil/strutil"
	"github.com/gookit/goutil/sysutil"
	"github.com/inhere/kite-go/pkg/util"
	"github.com/inhere/kite-go/pkg/xenv"
	"github.com/inhere/kite-go/pkg/xenv/shell"
)

var shellCmdOpts = struct {
	Type gflag.EnumString
	Install bool
	Reload  bool
	// Profile for pwsh. 无法区分版本，需要手动设置
	Profile string
}{
	Type: cflag.NewEnumString("bash", "zsh", "pwsh", "cmd"),
}

// ShellCmd the xenv shell command
var ShellCmd = &gcli.Command{
	Name: "shell",
	Desc: "Generate shell integration script contents",
	Help: `
<cyan>Auto Configure:</>
  # pwsh
  kite xenv shell --install -t pwsh --profile $PROFILE.CurrentUserAllHosts
  # bash, zsh
  kite xenv shell --install -t $SHELL

<cyan>Config for Bash:</>
  // write to .bashrc OR .bash_profile
  eval "$(kite xenv shell --type bash)"

<cyan>Config for Zsh:</>
  // write to .zshrc OR .zsh_profile
  eval "$(kite xenv shell --type zsh)"

<cyan>Config for Pwsh:</>
  # write expr to profile. (find by: echo $PROFILE.CurrentUserAllHosts)
  # Method 1:
  Invoke-Expression (&kite xenv shell --type pwsh)
  # Method 2:
  kite xenv shell --type pwsh | Out-String | Invoke-Expression
`,
	Config: func(c *gcli.Command) {
		c.BoolOpt(&shellCmdOpts.Reload, "reload", "r", false, "Reload the xenv shell script codes")
		c.BoolOpt(&shellCmdOpts.Install, "install", "i", false, "Install the xenv hook script to profile")
		c.VarOpt(&shellCmdOpts.Type, "type", "t", "Shell type (bash, zsh, pwsh, cmd)")
		c.StrOpt2(&shellCmdOpts.Profile, "profile", "current used PowerShell profile path.")
	},
	Func: func(c *gcli.Command, args []string) error {
		shellType, err := getShellType()
		if err != nil {
			return err
		}

		// Create env service
		envSvc, err := xenv.EnvService()
		if err != nil {
			return err
		}

		// 自动安装钩子脚本到 user shell 配置文件
		if shellCmdOpts.Install {
			return envSvc.WriteHookToProfile(shellType, shellCmdOpts.Profile)
		}

		// 生成钩子脚本
		hookScript, err := envSvc.GenHookScripts(shellType)
		if err != nil {
			return err
		}

		fmt.Print(hookScript)
		return nil
	},
}

func getShellType() (st shell.ShellType, err error) {
	var shellName string
	if shellCmdOpts.Reload {
		hookShellName := util.HookShell()
		if hookShellName == "" {
			return st, errorx.New("current is not in Shell hooking, XENV_HOOK_SHELL is empty")
		}
		shellName = hookShellName
	} else {
		shellName = strutil.OrElse(shellCmdOpts.Type.String(), sysutil.CurrentShell(true))
		if shellName == "" {
			return st, errorx.Errf("please specify the shell type (%s)", shellCmdOpts.Type.EnumString())
		}
	}

	shellType, err := shell.TypeFromString(shellName)
	if err != nil {
		return st, err
	}
	return shellType, nil
}

// HookInitCmd the xenv hook init command
//   - 通过配置上面的 shell 命令到 user 配置文件后，会自动执行该命令
//   - 调用当前命令，可以返回脚本内容自动执行
var HookInitCmd = &gcli.Command{
	Hidden: true, // This is an internal command
	Name:   "hook-init",
	Desc:   "Initialize the xenv hook script",
	Func: func(c *gcli.Command, args []string) error {
		return nil // TODO
	},
}
