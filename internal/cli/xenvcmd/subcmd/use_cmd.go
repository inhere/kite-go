package subcmd

import (
	"github.com/gookit/gcli/v3"
	"github.com/inhere/kite-go/pkg/xenv"
	"github.com/inhere/kite-go/pkg/xenv/shell"
)

var useCmdOpts = struct {
	Save bool
}{}

// UseCmd the xenv use command
var UseCmd = &gcli.Command{
	Name: "use",
	Help: "use [-g] <name:version>...",
	Desc: "Switch and activate different versions of development tools",
	Config: func(c *gcli.Command) {
		c.BoolOpt(&GlobalFlag, "global", "g", false, "Global operation, not the current session")
		c.BoolOpt(&useCmdOpts.Save, "save", "s", false, "Save the tools to current workdir .xenv.toml")

		c.AddArg("tools", "Name of the tool to activate, allow multi.", true, true)
	},
	Func: func(c *gcli.Command, args []string) error {
		// Create tool service
		toolSvc, err := xenv.ToolService()
		if err != nil {
			return err
		}

		useTools := c.Arg("tools").Strings()
		script, err1 := toolSvc.ActivateTools(useTools, GlobalFlag)
		if err1 == nil {
			shell.OutputScript(script)
		}
		return err1
	},
}

// UnuseCmd the xenv unuse command
var UnuseCmd = &gcli.Command{
	Name: "unuse",
	Help: "unuse [-g] <name:version>...",
	Desc: "Deactivate specific tool versions",
	Config: func(c *gcli.Command) {
		c.BoolOpt(&GlobalFlag, "global", "g", false, "Global operation, not the current session")
		c.AddArg("tools", "Name of the tool to deactivate, allow multi.", true, true)
	},
	Func: func(c *gcli.Command, args []string) error {
		// Create tool service
		toolSvc, err := xenv.ToolService()
		if err != nil {
			return err
		}

		unTools := c.Arg("tools").Strings()
		script, err1 := toolSvc.DeactivateTools(unTools, GlobalFlag)
		if err1 == nil {
			shell.OutputScript(script)
		}
		return err1
	},
}
