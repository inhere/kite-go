package subcmd

import (
	"github.com/gookit/gcli/v3"
	"github.com/inhere/kite-go/pkg/xenv"
	"github.com/inhere/kite-go/pkg/xenv/shell"
)

// NewUseCmd the xenv use command
func NewUseCmd() *gcli.Command {
	return &gcli.Command{
		Name: "use",
		Help: "use [-g] <name:version>...",
		Desc: "Switch and activate different versions of SDK/tool",
		Config: func(c *gcli.Command) {
			c.BoolOpt(&GlobalFlag, "global", "g", false, "Global operation, not the current session")
			c.BoolOpt(&SaveDirenv, "save", "s,d", false, "Save change to direnv config .xenv.toml")

			c.AddArg("tools", "Name of the tool to activate, allow multi.", true, true)
		},
		Func: func(c *gcli.Command, args []string) error {
			// Create tool service
			toolSvc, err := xenv.ToolService()
			if err != nil {
				return err
			}

			useTools := c.Arg("tools").Strings()
			script, err1 := toolSvc.ActivateSDKs(useTools, GetOpFlag(SaveDirenv, GlobalFlag))
			if err1 == nil {
				shell.OutputScript(script)
			}
			return err1
		},
	}
}

// NewUnuseCmd the xenv unuse command
func NewUnuseCmd() *gcli.Command {
	return &gcli.Command{
		Name: "unuse",
		Help: "unuse [-g] <name:version>...",
		Desc: "Deactivate specific SDK/tool versions",
		Config: func(c *gcli.Command) {
			c.BoolOpt(&GlobalFlag, "global", "g", false, "Global operation, not the current session")
			c.BoolOpt(&SaveDirenv, "save", "s,d", false, "Save change to direnv config .xenv.toml")
			c.AddArg("tools", "Name of the tool to deactivate, allow multi.", true, true)
		},
		Func: func(c *gcli.Command, args []string) error {
			// Create tool service
			toolSvc, err := xenv.ToolService()
			if err != nil {
				return err
			}

			unTools := c.Arg("tools").Strings()
			script, err1 := toolSvc.DeactivateSDKs(unTools, GetOpFlag(SaveDirenv, GlobalFlag))
			if err1 == nil {
				shell.OutputScript(script)
			}
			return err1
		},
	}
}
