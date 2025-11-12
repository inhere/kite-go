package subcmd

import (
	"fmt"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/gflag"
	"github.com/gookit/goutil/cflag"
	"github.com/gookit/goutil/errorx"
	"github.com/gookit/goutil/strutil"
	"github.com/gookit/goutil/sysutil"
	"github.com/inhere/kite-go/pkg/xenv"
	"github.com/inhere/kite-go/pkg/xenv/shell"
	"github.com/inhere/kite-go/pkg/xenv/xenvcom"
)

type shellOptions = struct {
	Type gflag.EnumString
	Install bool
	Reload  bool
	// Profile for pwsh. 无法区分版本，需要手动设置
	Profile string
}

// NewShellCmd the xenv shell command
func NewShellCmd() *gcli.Command {
	var shellOpts = shellOptions{
		Type: cflag.NewEnumString(shell.AllTypeStrings...),
	}

	return &gcli.Command{
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
			c.BoolOpt(&shellOpts.Reload, "reload", "r", false, "Reload the xenv shell script codes")
			c.BoolOpt(&shellOpts.Install, "install", "i", false, "Install the xenv hook script to profile")
			c.VarOpt(&shellOpts.Type, "type", "t", "Shell type (bash, zsh, pwsh, cmd)")
			c.StrOpt2(&shellOpts.Profile, "profile", "current used PowerShell profile path.")
		},
		Func: func(c *gcli.Command, args []string) error {
			shellType, err := getShellType(&shellOpts)
			if err != nil {
				return err
			}

			// Create env service
			toolSvc, err := xenv.ToolService()
			if err != nil {
				return err
			}

			// 自动安装钩子脚本到 user shell 配置文件
			if shellOpts.Install {
				return toolSvc.WriteHookToProfile(shellType, shellOpts.Profile)
			}

			// 生成钩子脚本
			hookScript, err := toolSvc.GenHookScripts(shellType)
			if err != nil {
				return err
			}

			fmt.Print(hookScript)
			return nil
		},
	}
}

func getShellType(shellOpts *shellOptions) (st shell.ShType, err error) {
	var shellName string
	if shellOpts.Reload {
		hookShellName := xenvcom.HookShell()
		if hookShellName == "" {
			return st, errorx.New("current is not in Shell hooking, XENV_HOOK_SHELL is empty")
		}
		shellName = hookShellName
	} else {
		shellName = strutil.OrElse(shellOpts.Type.String(), sysutil.CurrentShell(true))
		if shellName == "" {
			return st, errorx.Errf("please specify the shell type (%s)", shellOpts.Type.EnumString())
		}
	}

	shellType, err := shell.TypeFromString(shellName)
	if err != nil {
		return st, err
	}
	return shellType, nil
}

// ShellHookInitCmd the xenv hook init command
//   - 配置了 xenv shell 命令到 user 配置文件后，会自动执行该命令
//   - 调用当前命令，可以返回脚本内容自动执行
func ShellHookInitCmd() *gcli.Command {
	var initHookOpts = struct {
		Type gflag.EnumString
	}{
		Type: cflag.NewEnumString(shell.AllTypeStrings...),
	}

	return &gcli.Command{
		Name:   "shell-init-hook",
		Desc:   "Initialize the xenv hook script",
		Hidden: true, // This is an internal command
		Config: func(c *gcli.Command) {
			c.VarOpt(&initHookOpts.Type, "type", "t", "Shell type (bash, zsh, pwsh, cmd)")
		},
		Func: func(c *gcli.Command, args []string) error {
			return nil // TODO
		},
	}
}

// ShellDirenvCmd the xenv init shell direnv command
//  - 仅在配置了 xenv shell 命令时，cd 到新目录会自动调用当前命令
//  - 监听进入目录时，自动检测 .xenv.toml 文件，并加载里面的配置
func ShellDirenvCmd() *gcli.Command {
	var direnvOpts = struct {
		Type gflag.EnumString
	}{
		Type: cflag.NewEnumString(shell.AllTypeStrings...),
	}

	// into-dir, leave-dir
	return &gcli.Command{
		Name:    "shell-direnv",
		Desc:    "Initialize direnv state on current workdir",
		Hidden:  true, // This is an internal command
		Aliases: []string{"init-direnv", "into-dir"},
		Config: func(c *gcli.Command) {
			c.VarOpt(&direnvOpts.Type, "type", "t", "Shell type (bash, zsh, pwsh, cmd)")
		},
		Func: func(c *gcli.Command, args []string) error {
			// Create tool service
			toolSvc, err := xenv.ToolService()
			if err != nil {
				return err
			}

			script, err1 := toolSvc.SetupDirenv()
			if err1 == nil {
				shell.OutputScript(script)
			}
			return err1
		},
	}
}
