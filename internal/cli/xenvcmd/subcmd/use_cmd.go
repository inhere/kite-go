package subcmd

import (
	"github.com/gookit/gcli/v3"
	"github.com/inhere/kite-go/pkg/xenv"
	"github.com/inhere/kite-go/pkg/xenv/models"
	"github.com/inhere/kite-go/pkg/xenv/shell"
)

// NewUseCmd the xenv use command
func NewUseCmd() *gcli.Command {
	var useCmdOpts = struct {
		Save bool
	}{}

	return &gcli.Command{
		Name: "use",
		Help: "use [-g] <name:version>...",
		Desc: "Switch and activate different versions of SDK/tool",
		Config: func(c *gcli.Command) {
			c.BoolOpt(&GlobalFlag, "global", "g", false, "Global operation, not the current session")
			c.BoolOpt(&useCmdOpts.Save, "save", "s", false, "Save change to current workdir .xenv.toml")

			c.AddArg("tools", "Name of the tool to activate, allow multi.", true, true)
		},
		Func: func(c *gcli.Command, args []string) error {
			// Create tool service
			toolSvc, err := xenv.ToolService()
			if err != nil {
				return err
			}

			useTools := c.Arg("tools").Strings()
			script, err1 := toolSvc.ActivateSDKs(useTools, models.GetOpFlag(useCmdOpts.Save, GlobalFlag))
			if err1 == nil {
				shell.OutputScript(script)
			}
			return err1
		},
	}
}

// NewUnuseCmd the xenv unuse command
func NewUnuseCmd() *gcli.Command {
	var unuseCmdOpts = struct {
		Save bool
	}{}

	return &gcli.Command{
		Name: "unuse",
		Help: "unuse [-g] <name:version>...",
		Desc: "Deactivate specific SDK/tool versions",
		Config: func(c *gcli.Command) {
			c.BoolOpt(&GlobalFlag, "global", "g", false, "Global operation, not the current session")
			c.BoolOpt(&unuseCmdOpts.Save, "save", "s", false, "Save change to current workdir .xenv.toml")
			c.AddArg("tools", "Name of the tool to deactivate, allow multi.", true, true)
		},
		Func: func(c *gcli.Command, args []string) error {
			// Create tool service
			toolSvc, err := xenv.ToolService()
			if err != nil {
				return err
			}

			unTools := c.Arg("tools").Strings()
			script, err1 := toolSvc.DeactivateSDKs(unTools, models.GetOpFlag(unuseCmdOpts.Save, GlobalFlag))
			if err1 == nil {
				shell.OutputScript(script)
			}
			return err1
		},
	}
}
