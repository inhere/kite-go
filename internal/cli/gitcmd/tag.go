package gitcmd

import (
	"errors"

	"github.com/gookit/color/colorp"
	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/gflag"
	"github.com/gookit/gcli/v3/interact"
	"github.com/gookit/gcli/v3/show"
	"github.com/gookit/gitw"
	"github.com/gookit/gitw/gitutil"
)

// TagCmd instance
var TagCmd = &gcli.Command{
	Name: "tag",
	Desc: "extra git tag commands",
	Subs: []*gcli.Command{
		TagListCmd,
		TagCreateCmd,
		TagDeleteCmd,
	},
}

// TagListCmd instance
var TagListCmd = &gcli.Command{
	Name:    "list",
	Aliases: []string{"ls"},
	Desc:    "list tags for the git repository",
	Func: func(c *gcli.Command, args []string) error {
		return errors.New("TODO")
	},
}

var tcOpts = struct {
	Next    bool   `flag:"create next version tag;false;;n"`
	Hash    string `flag:"create tag by commit hash;false;;cid"`
	Message string `flag:"tag message;false;;m"`
	Version string `flag:"tag version, eg: v2.0.1;false;;v"`
}{}

// TagCreateCmd instance
var TagCreateCmd = &gcli.Command{
	Name:    "create",
	Aliases: []string{"new"},
	Desc:    "create new tag by `git tag`",
	Help: `
# Examples:
  {$fullCmd} --next
  {$fullCmd} -v v2.0.1
`,
	Config: func(c *gcli.Command) {
		c.MustFromStruct(&tcOpts, gflag.TagRuleSimple)
	},
	Func: func(c *gcli.Command, args []string) error {
		lp := gitw.NewRepo(GitOpts.Workdir).
			PrintCmdOnExec().
			SetDryRun(GitOpts.DryRun)

		err := lp.Cmd("pull", "-np").Run()
		if err != nil {
			return err
		}

		// git fetch --tags
		err = lp.Cmd("fetch", "--tags").Run()
		if err != nil {
			return err
		}

		lastVer := lp.LargestTag()
		nextVer := tcOpts.Version
		if len(nextVer) == 0 {
			nextVer = gitutil.NextVersion(lastVer)
		} else {
			var ok bool
			nextVer, ok = gitutil.FormatVersion(nextVer)
			if !ok {
				return c.NewErrf("invalid version: %s", nextVer)
			}
		}

		nextVer = "v" + nextVer
		message := tcOpts.Message
		if len(message) == 0 {
			message = ":bookmark: release the " + nextVer
		}

		show.AList("create tag", map[string]any{
			"Hash ID":     tcOpts.Hash,
			"Prev tag":    lastVer,
			"New tag":     nextVer,
			"Tag Message": message,
			"Dry Run":     GitOpts.DryRun,
		})

		if interact.Unconfirmed("Ensure create and push new tag?", true) {
			colorp.Infoln("Quit, Bye!")
			return nil
		}

		err = lp.Cmd("tag", "-a", nextVer, "-m", message).Run()
		if err != nil {
			return err
		}

		err = lp.Cmd("push", "origin", nextVer).Run()
		if err != nil {
			return err
		}

		colorp.Successf("Successful create tag: %s\n", nextVer)
		return nil
	},
}

// TagDeleteCmd instance
var TagDeleteCmd = &gcli.Command{
	Name:    "delete",
	Aliases: []string{"del", "rm", "remove"},
	Desc:    "delete exists tags by `git tag`",
	Func: func(c *gcli.Command, args []string) error {
		return errors.New("TODO")
	},
}
