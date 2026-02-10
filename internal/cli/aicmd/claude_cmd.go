package aicmd

import (
	"fmt"

	"github.com/gookit/gcli/v3"
	"github.com/inhere/kite-go/internal/service"
	"github.com/inhere/kite-go/internal/service/aiservice"
)

// ClaudeCommand is the AI Claude command group
var ClaudeCommand = &gcli.Command{
	Name: "claude",
	Desc: "Claude code managment commands",
	Aliases: []string{"cc"},
	Subs: []*gcli.Command{
		NewClaudeApiCmd(),
	},
}

// NewClaudeApiCmd creates the Claude API info manage command
func NewClaudeApiCmd() *gcli.Command {
	var opts = aiservice.SetAuthInfoParam{}

	return &gcli.Command{
		Name: "api",
		Desc: "Claude API config manage operations",
		Config: func(c *gcli.Command) {
			c.StrOpt2(&opts.Provider, "use", "use the model provider, allowed: glm,minimax,kimi,claude")
			c.StrOpt2(&opts.Shell, "shell", "the shell type, allowed: bash,pwsh")
			c.BoolOpt2(&opts.Write, "write,w", "whether to write the config to file")
			c.BoolOpt2(&opts.Show, "show,i", "show the config information")
		},
		Help: `
# active in bash, pwsh shell
eval "$({$binWithCmd} --shell bash --use kimi)"
# active in pwsh shell
{$binWithCmd} --shell pwsh --use kimi | Out-String | Invoke-Expression
`,
		Func: func(cmd *gcli.Command, args []string) error {
			aisrv, err := service.AI().Init()
			if err != nil {
				return fmt.Errorf("failed to initialize AI service: %w", err)
			}
			return aisrv.SetAuthInfo(opts)
		},
	}
}
