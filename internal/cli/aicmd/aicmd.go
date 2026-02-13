package aicmd

import (
	"fmt"

	"github.com/gookit/gcli/v3"
	"github.com/inhere/kite-go/internal/service"
)

var aiOpts = struct {
	showConfig bool
	configKey string
}{}

var AICommand = &gcli.Command{
	Name: "ai",
	Desc: "AI tool command",
	Subs: []*gcli.Command{
		NewAIChatCmd(),
		NewQuestionCmd(),
		NewTranslateCmd(),
		NewAIToolCmd(),
		ClaudeCommand,
	},
	Config: func(c *gcli.Command) {
		c.BoolOpt2(&aiOpts.showConfig, "show", "Show config info")
		c.StrOpt2(&aiOpts.configKey, "key", "Show config info by key, with --show")
	},
	Func: func(c *gcli.Command, args []string) error {
		if aiOpts.showConfig {
			aisrv, err := service.AI().Init()
			if err != nil {
				return fmt.Errorf("failed to initialize AI service: %w", err)
			}
			aisrv.ShowConfig(aiOpts.configKey)
			return nil
		}
		return c.ShowHelp()
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
