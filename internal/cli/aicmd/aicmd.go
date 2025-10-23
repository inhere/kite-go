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
		NewAIToolCmd(),
	},
}

type AICommonOptions struct {
	// custom set llm model name
	Model string
	// custom set llm provider name. eg: openai, deepseek, aliyun/bailian, siliconflow
	Provider string
	SystemMsg string
}

// BindFlags for common options
func (o *AICommonOptions) BindFlags(c *gcli.Command) {
	c.StrOpt2(&o.Model, "model,m", "custom set llm model name. default: DEFAULT_LLM_MODEL")
	c.StrOpt2(&o.Provider, "provider,p", `custom set llm provider name.
eg: openai, ds/deepseek, ali/aliyun/bailian, sf/siliconflow, kimi, zai/zhipu`)
	c.StrOpt2(&o.SystemMsg, "system,s", "custom set system prompt message")
}
