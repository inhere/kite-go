package aicmd

import (
	"github.com/gookit/gcli/v3"
)

var AICommand = &gcli.Command{
	Name: "ai",
	Desc: "AI tool command",
	Subs: []*gcli.Command{
		NewAIChatCmd(),
		NewQuestionCmd(),
		NewTranslateCmd(),
	},
}

type AICommonOptions struct {
	// custom set llm model name
	Model string
	// custom set llm provider name. eg: openai, deepseek, aliyun/bailian, siliconflow
	Provider string
}

// BindFlags for common options
func (o *AICommonOptions) BindFlags(c *gcli.Command) {
	c.StrOpt(&o.Model, "model", "", "")
	c.StrOpt(&o.Provider, "provider", "", "")
}
