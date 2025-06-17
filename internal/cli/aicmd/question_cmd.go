package aicmd

import "github.com/gookit/gcli/v3"

// NewQuestionCmd create.
func NewQuestionCmd() *gcli.Command {
	var askOpts = struct {
		AICommonOptions
		interactive bool
	}{}

	return &gcli.Command{
		Name:    "question",
		Aliases: []string{"q", "ask"},
		Desc:    "Ask the AI questions and get the results",
		Config: func(c *gcli.Command) {
			askOpts.BindFlags(c)
			c.BoolOpt2(&askOpts.interactive, "interactive,i", "into interactive mode")
		},
		Func: func(c *gcli.Command, args []string) error {
			return nil
		},
	}
}
