package subcmd

import (
	"fmt"

	"github.com/gookit/gcli/v3"
	"github.com/inhere/kite-go/pkg/xenv/config"
	"github.com/inhere/kite-go/pkg/xenv/tools"
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
		useTools := c.Arg("tools").Strings()
		for _, arg := range useTools {
			// Parse name:version
			name, version, err := parseNameVersion(arg)
			if err != nil {
				return err
			}

			// Initialize configuration
			if err := config.Mgr.Init(); err != nil {
				return fmt.Errorf("failed to initialize configuration: %w", err)
			}

			// Create activator
			activator := tools.NewActivator(config.Mgr.Config, config.Mgr.State)

			// Activate the tool
			if err := activator.ActivateTool(name, version, GlobalFlag); err != nil {
				return fmt.Errorf("failed to activate tool %s:%s: %w", name, version, err)
			}

			// Save configuration if global flag is set
			if GlobalFlag {
				c.Infof("Set %s:%s as global default\n", name, version)
			} else {
				c.Infof("Set %s:%s for current session\n", name, version)
			}
		}

		return nil
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
		unTools := c.Arg("tools").Strings()

		for _, arg := range unTools {
			// Parse name:version
			name, version, err := parseNameVersion(arg)
			if err != nil {
				return err
			}

			// Initialize configuration
			if err := config.Mgr.Init(); err != nil {
				return fmt.Errorf("failed to initialize configuration: %w", err)
			}

			// Create activator
			activator := tools.NewActivator(config.Mgr.Config, config.Mgr.State)

			// Deactivate the tool
			if err := activator.DeactivateTool(name, version, GlobalFlag); err != nil {
				return fmt.Errorf("failed to deactivate tool %s:%s: %w", name, version, err)
			}

			// Save configuration if global flag is set
			if GlobalFlag {
				c.Infof("Unset %s:%s from global default\n", name, version)
			} else {
				c.Infof("Unset %s:%s from current session\n", name, version)
			}
		}

		return nil
	},
}
