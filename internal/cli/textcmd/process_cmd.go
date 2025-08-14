package textcmd

import "github.com/gookit/gcli/v3"

// NewProcessCmd create a new ProcessCmd instance
// 实现处理输入文本内容
func NewProcessCmd() *gcli.Command {
	var procOpts = struct {
	}{}

	return &gcli.Command{
		Name:    "process",
		Desc:    "Process input text contents",
		Aliases: []string{"proc", "handle"},
		Config: func(c *gcli.Command) {
			c.MustFromStruct(&procOpts)
		},
		Func: func(c *gcli.Command, args []string) error {
			// TODO
			return nil
		},
	}
}

// TODO 实现简单的 set 命令，
//  - 可以根据匹配行，更新行内容
//  - 可以根据匹配行，追加新内容
// func NewSetCmd()
