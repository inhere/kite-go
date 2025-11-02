package subcmd

import (
	"fmt"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/x/ccolor"
	"github.com/inhere/kite-go/pkg/xenv"
)

var (
	// GlobalFlag option value
	GlobalFlag bool
	SaveDirenv bool
	DebugMode  bool
)

// EnvCmd the xenv env command
var EnvCmd = &gcli.Command{
	Name: "env",
	Desc: "Manage environment variables",
	Subs: []*gcli.Command{
		EnvSetCmd(),
		EnvUnsetCmd(),
		EnvListCmd(),
	},
	Func: func(c *gcli.Command, args []string) error {
		return listEnvs()
	},
}

// EnvSetCmd command for setting environment variables
//
// Test run:
//
//	// pwsh
//	$env:XENV_HOOK_SHELL="pwsh"; kite xenv set TEST003 value003
func EnvSetCmd() *gcli.Command {
	return &gcli.Command{
		Name: "set",
		Help: "set [-g] [-s|-d] <name> <value>",
		Desc: "Set an environment variable",
		Config: func(c *gcli.Command) {
			c.BoolOpt(&SaveDirenv, "direnv", "s,d", false, "Save change to direnv config .xenv.toml")
			c.BoolOpt(&GlobalFlag, "global", "g", false, "Operate for global config")

			c.AddArg("name", "environment key name", true)
			c.AddArg("value", "environment value", true)
		},
		Func: func(c *gcli.Command, args []string) error {
			name := c.Arg("name").String()
			value := c.Arg("value").String()

			// Create env service
			envSvc, err := xenv.EnvService()
			if err != nil {
				return err
			}

			// Set the environment variable
			script, err := envSvc.SetEnv(name, value, GetOpFlag())
			if err != nil {
				return fmt.Errorf("failed to set environment variable: %w", err)
			}

			// Save configuration if global
			if GlobalFlag {
				ccolor.Infof("Set %s=%s globally\n", name, value)
			} else {
				ccolor.Infof("Set %s=%s for current session\n", name, value)
			}

			if script != "" {
				fmt.Printf("%s\n%s\n", xenv.ScriptMark, script)
			}
			return nil
		},
	}
}

// EnvUnsetCmd command for unsetting environment variables
func EnvUnsetCmd() *gcli.Command {
	return &gcli.Command{
		Name: "unset",
		Help: "unset [-g] [-s|-d] <name...>",
		Desc: "Unset environment variables",
		Config: func(c *gcli.Command) {
			c.BoolOpt(&SaveDirenv, "direnv", "s,d", false, "Operate for direnv config .xenv.toml")
			c.BoolOpt(&GlobalFlag, "global", "g", false, "Operate for global config")
			c.AddArg("names", "environment key name", true, true)
		},
		Func: func(c *gcli.Command, args []string) error {
			// Create env service
			envSvc, err := xenv.EnvService()
			if err != nil {
				return err
			}

			names := c.Arg("names").Strings()

			// Unset the environment variables
			script, err1 := envSvc.UnsetEnvs(names, GetOpFlag())
			if err1 != nil {
				return fmt.Errorf("failed to set environment variable: %w", err1)
			}

			// Save configuration if global
			if GlobalFlag {
				ccolor.Infof("Unset %s globally\n", names)
			} else {
				ccolor.Infof("Unset %s for current session\n", names)
			}

			if script != "" {
				fmt.Printf("%s\n%s\n", xenv.ScriptMark, script)
			}
			return nil
		},
	}
}

// EnvListCmd command for listing environment variables
func EnvListCmd() *gcli.Command {
	return &gcli.Command{
		Name:    "list",
		Desc: "List environment variables",
		Aliases: []string{"ls"},
		Func: func(c *gcli.Command, args []string) error {
			return listEnvs()
		},
	}
}

func listEnvs() error {
	// Create env service
	envSvc, err := xenv.EnvService()
	if err != nil {
		return err
	}

	// List environment variables
	envVars := envSvc.GlobalEnv()
	ccolor.Infoln("Global Environment Variables:")
	for name, envVar := range envVars {
		fmt.Printf("  %s=%s\n", name, envVar)
	}

	if envSvc.IsSessionEnv() {
		sessVars := envSvc.SessionEnv()
		ccolor.Infoln("Session Environment Variables:")
		for name, envVar := range sessVars {
			fmt.Printf("  %s=%s\n", name, envVar)
		}
	}
	return nil
}
