package envcmd

import "github.com/gookit/gcli/v3"

// NewEnvManageCmd 创建环境管理命令
// 参考 mise 等工具
func NewEnvManageCmd() *gcli.Command {
	return &gcli.Command{
		Name:    "env",
		Aliases: []string{"envs"},
		Desc:    "manage local development environment SDK",
		Help: `
Commands Usage:
  use <sdk:version>...     Activate SDK versions
    -s, --save             Save configuration to project file
  unuse <sdk>...           Deactivate SDKs
  add <sdk:version>...     Download and install SDK versions
  list [sdk]               List installed SDKs

Examples:
  ktenv use node:18 go:1.21
  ktenv use -s node:lts
  ktenv unuse node
  ktenv add go:1.22
  ktenv list
  ktenv list go

Supported SDKs:
  go, node, java, flutter

Version formats:
  <sdk>:<version>         Exact version (go:1.21.5)
  <sdk>:<major>           Latest patch version (node:18)
  <sdk>:lts               Long-term support version
  <sdk>:latest            Latest stable version
`,
		Subs: []*gcli.Command{
			NewEnvListCmd(),
			NewEnvAddCmd(),
			NewEnvRemoveCmd(),
			NewEnvUseCmd(),
			NewEnvUnuseCmd(),
			NewEnvShellCmd(),
			NewEnvConfigCmd(),
		},
	}
}

func NewEnvListCmd() *gcli.Command {
	return &gcli.Command{
		Name:    "list",
		Desc:    "list local environment SDK",
		Aliases: []string{"ls", "l"},
		Config: func(c *gcli.Command) {
			c.StrOpt(&envListOpts.sdkType, "type", "t", "filter by SDK type")
		},
		Func: handleEnvList,
	}
}

var envListOpts = struct {
	sdkType string
}{}

// NewEnvShellCmd 创建环境shell注入命令
// 将会创建环境shell注入脚本代码，并输出到标准输出。
//
// 内部会生成一个 ktenv shell 函数，用户通过 ktenv 函数实现切换shell环境信息。
//
// Usage: kite dev env shell [pwsh|bash|cmd|zsh]
func NewEnvShellCmd() *gcli.Command {
	return &gcli.Command{
		Name:    "shell",
		Desc:    "create environment shell injection scripts",
		Aliases: []string{"sh"},
		Func:    handleEnvShell,
	}
}
