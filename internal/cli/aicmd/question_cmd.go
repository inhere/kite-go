package aicmd

import (
	"context"
	"fmt"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/envutil"
	"github.com/gookit/goutil/strutil"
	"github.com/inhere/kite-go/internal/apputil"
	"github.com/sashabaranov/go-openai"
)

// NewQuestionCmd create.
func NewQuestionCmd() *gcli.Command {
	var askOpts = struct {
		AICommonOptions
		interactive bool
		// 内置的角色名称 go-dev, java-dev, android-dev 等，选择后将会使用对应的提示词设置系统提示词
		roleName string
	}{}

	return &gcli.Command{
		Name:    "question",
		Aliases: []string{"q", "ask"},
		Desc:    "Ask the AI questions and get the results",
		Config: func(c *gcli.Command) {
			askOpts.BindFlags(c)
			c.BoolOpt2(&askOpts.interactive, "interactive,i", "into interactive mode")

			c.AddArg("question", "The question to ask", true)
		},
		Func: func(c *gcli.Command, args []string) error {
			question, err := apputil.ReadSource(c.Arg("question").String())
			if err != nil {
				return err
			}

			client := apputil.NewOpenaiClient()
			modelName := envutil.GetOne([]string{"QUESTION_MODEL", "DEFAULT_LLM_MODEL"}, "gpt-3.5-turbo")
			messages := []openai.ChatCompletionMessage{
				// system message
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: strutil.OrElse(askOpts.SystemMsg, "You are a helpful assistant."),
				},
				// user message
				{
					Role:    openai.ChatMessageRoleUser,
					Content: question,
				},
			}

			c.Infof("Questioning AI Model(%s) ...\n", modelName)
			resp, err := client.CreateChatCompletion(
				context.Background(),
				openai.ChatCompletionRequest{
					// Model: openai.GPT3Dot5Turbo,
					Model:    modelName,
					Messages: messages,
				},
			)

			if err == nil {
				fmt.Println("ANSWER:\n", resp.Choices[0].Message.Content)
			}
			return err
		},
	}
}
