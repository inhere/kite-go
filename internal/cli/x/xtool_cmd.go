package x

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/errorx"
)

// XToolCmd 管理安装本机需要使用的工具(主要是命令工具 eg: fzf, git 等等)
var XToolCmd = &gcli.Command{
	Name: "xtool",
	Desc: "Unified installation and management of external tools",
	Subs: []*gcli.Command{
	},
	Func: func(c *gcli.Command, args []string) error {
		return errorx.Raw("TODO") // TODO
	},
}
