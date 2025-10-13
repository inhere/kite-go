package envcmd

import (
	"fmt"

	"github.com/gookit/gcli/v3"
	"github.com/inhere/kite-go/pkg/envmgr"
)

var envUseOpts = struct {
	save bool
}{}

func NewEnvUseCmd() *gcli.Command {
	return &gcli.Command{
		Name:    "use",
		Desc:    "enable new tool or version environment",
		Aliases: []string{"switch", "set"},
		Config: func(c *gcli.Command) {
			c.BoolOpt(&envUseOpts.save, "save", "s", false, "save configuration to project file")

			c.AddArg("sdk_versions", "specify sdk and version to use. <sdk:version> [sdk:version...]", true, true)
		},
		Func: handleEnvUse,
	}
}

// handleEnvUse 处理环境使用命令
func handleEnvUse(c *gcli.Command, _ []string) error {
	specs := c.Arg("sdk_versions").Strings()
	em := getEnvManager()

	for _, specStr := range specs {
		spec, err := envmgr.ParseVersionSpec(specStr)
		if err != nil {
			return fmt.Errorf("invalid version spec '%s': %w", specStr, err)
		}

		if err := em.UseSDK(spec.SDK, spec.Version, envUseOpts.save); err != nil {
			return fmt.Errorf("failed to use %s: %w", specStr, err)
		}

		c.Printf("Activated %s:%s\n", spec.SDK, spec.Version)
	}

	return nil
}

// NewEnvUnuseCmd 处理环境取消使用命令
func NewEnvUnuseCmd() *gcli.Command {
	return &gcli.Command{
		Name:    "unuse",
		Desc:    "unuse environment",
		Aliases: []string{"deactivate", "unset"},
		Func: func(c *gcli.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("usage: kite dev env unuse <sdk>")
			}

			em := getEnvManager()

			for _, sdk := range args {
				if err := em.UnuseSDK(sdk); err != nil {}
			}
			return nil
		},
	}
}
