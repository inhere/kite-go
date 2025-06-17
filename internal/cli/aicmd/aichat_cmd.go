package aicmd

import (
	"os"

	"github.com/gookit/gcli/v3"
	"golang.org/x/term"
)

// NewAIChatCmd create 实现一个AI终端
//   - 可以跟AI大模型交互
//   - 支持多轮对话
//   - 支持子命令调整设置
//   - 支持上下文存储
//   - 支持历史记录存储
func NewAIChatCmd() *gcli.Command {
	var aiTerm = AITerminal{}

	return &gcli.Command{
		Name: "chat",
		Desc: "Chat with AI",
		Config: func(c *gcli.Command) {
			aiTerm.BindFlags(c)
		},
		Func: func(c *gcli.Command, args []string) error {
			aiTerm.Create()

			aiTerm.Run()
			return nil
		},
	}
}

type AITerminal struct {
	AICommonOptions
	SystemPrompt string

	term *term.Terminal
}

// Create terminal
func (at *AITerminal) Create() {
	at.term = term.NewTerminal(os.Stdout, "> ")
}

func (at *AITerminal) Run() {

}
