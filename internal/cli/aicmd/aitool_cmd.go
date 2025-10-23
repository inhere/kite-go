package aicmd

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/errorx"
)

// NewAIToolCmd AI工具命令,通过内置的提示词将ai当做工具使用
func NewAIToolCmd() *gcli.Command {
	return &gcli.Command{
		Name: "tool",
		Desc: "Use AI as tools do something. such as convert, format and more",
		Help: `
Available AI Tools:
 <cyan>Convert Format</>:
  json-to-yaml
  json-to-toml
  yaml-to-json
  yaml-to-toml
 <cyan>Convert Language</>:
  java-to-kotlin
  kotlin-to-java
`,
		Func: func(c *gcli.Command, args []string) error {
			// TODO
			return errorx.Raw("Not implemented yet")
		},
	}
}
