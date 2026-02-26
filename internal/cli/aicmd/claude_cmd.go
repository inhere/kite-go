package aicmd

import (
	"fmt"

	"github.com/gookit/gcli/v3"
	"github.com/inhere/kite-go/internal/service"
	"github.com/inhere/kite-go/internal/service/aiservice"
)

// ClaudeCommand is the AI Claude command group
var ClaudeCommand = &gcli.Command{
	Name:    "claude",
	Desc:    "Claude code managment commands",
	Aliases: []string{"cc"},
	Subs: []*gcli.Command{
		NewClaudeSetCmd(),
	},
}

// NewClaudeSetCmd creates the Claude Code info manage command
func NewClaudeSetCmd() *gcli.Command {
	var opts = aiservice.SetClaudeCodeParam{}

	return &gcli.Command{
		Name:    "set",
		Aliases: []string{"api"},
		Desc:    "Claude code config manage operations",
		Config: func(c *gcli.Command) {
			c.StrOpt2(&opts.Provider, "use", "use the model provider, allowed: glm,minimax,kimi,claude")
			c.StrOpt(&opts.KeyName, "key", "", "default", "the api key name in api_keys")
			c.StrOpt(&opts.Model, "model", "m", "", "model name for use, default from config")
			c.StrOpt(&opts.Scope, "scope", "s", "user", "scope on write config, allow: user, project")
			c.StrOpt2(&opts.Shell, "shell", "the shell type, allowed: bash,pwsh")
			c.BoolOpt2(&opts.Write, "write,w", "whether to write the config to file")
			c.BoolOpt2(&opts.Show, "show,i", "show the config information")
			c.BoolOpt2(&opts.List, "list,l", "list all api config information")
		},
		Help: `
# active in bash, zsh shell
eval "$({$binWithCmd} --shell bash --use kimi)"
# active in pwsh shell
{$binWithCmd} --shell pwsh --use kimi | Out-String | Invoke-Expression
`,
		Func: func(cmd *gcli.Command, args []string) error {
			aisrv, err := service.AI().Init()
			if err != nil {
				return fmt.Errorf("failed to initialize AI service: %w", err)
			}
			return aisrv.SetClaudeCode(opts)
		},
	}
}
