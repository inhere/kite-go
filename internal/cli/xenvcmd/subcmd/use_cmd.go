package subcmd

import (
	"github.com/gookit/gcli/v3"
	"github.com/inhere/kite-go/pkg/xenv"
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
		c.BoolOpt(&useCmdOpts.Save, "save", "s", false, "Save the tool version to current workdir .xenv.toml")

		c.AddArg("tools", "Name of the tool to install, allow multi.", true, true)
	},
	Func: func(c *gcli.Command, args []string) error {
		// Create tool service
		toolSvc, err := xenv.ToolService()
		if err != nil {
			return err
		}

		useTools := c.Arg("tools").Strings()
		return toolSvc.ActivateTools(useTools, GlobalFlag)
	},
}

// UnuseCmd the xenv unuse command
var UnuseCmd = &gcli.Command{
	Name: "unuse",
	Help: "unuse [-g] <name:version>...",
	Desc: "Deactivate specific tool versions",
	Config: func(c *gcli.Command) {
		c.BoolOpt(&GlobalFlag, "global", "g", false, "Global operation, not the current session")
		c.AddArg("tools", "Name of the tool to install, allow multi.", true, true)
	},
	Func: func(c *gcli.Command, args []string) error {
		// Create tool service
		toolSvc, err := xenv.ToolService()
		if err != nil {
			return err
		}

		unTools := c.Arg("tools").Strings()
		return toolSvc.DeactivateTools(unTools, GlobalFlag)
	},
}
