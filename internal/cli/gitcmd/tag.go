package gitcmd

import (
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
		lp := gitw.NewRepo(GitOpts.Workdir).
			PrintCmdOnExec().
			SetDryRun(GitOpts.DryRun)

		// 获取标签列表
		tags := lp.Tags()
		if lp.Err() != nil {
			return lp.Err()
		}

		if len(tags) == 0 {
			colorp.Infof("No tags found in the repository.\n")
			return nil
		}

		// 打印标签列表
		colorp.Infof("Tags in the repository:\n")
		for _, tag := range tags {
			colorp.Cyanf("  %s\n", tag)
		}

		return nil
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
	Func: func(c *gcli.Command, _ []string) error {
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
	Config: func(c *gcli.Command) {
		c.AddArg("tags", "tags name to delete", true, true)
	},
	Func: func(c *gcli.Command, _ []string) error {
		tagNames := c.Arg("tags").Strings()
		if len(tagNames) == 0 {
			return c.NewErr("please provide at least one tag name to delete")
		}

		repo := gitw.NewRepo(GitOpts.Workdir).
			PrintCmdOnExec().
			SetDryRun(GitOpts.DryRun)

		// 确认删除操作
		colorp.Infof("Tags to delete: %v\n", tagNames)
		if interact.Unconfirmed("Are you sure you want to delete these tags?", true) {
			colorp.Infoln("Operation cancelled.")
			return nil
		}

		// 删除本地标签
		colorp.Cyanf("Deleting local tags...\n")
		for _, tag := range tagNames {
			err := repo.Cmd("tag", "-d", tag).Run()
			if err != nil {
				return err
			}
			colorp.Successf("Deleted local tag: %s\n", tag)
		}

		// 删除远程标签
		colorp.Cyanf("Deleting remote tags...\n")
		argsWithPrefix := make([]string, len(tagNames))
		for i, tag := range tagNames {
			argsWithPrefix[i] = ":" + tag
		}

		// RUN: git push origin :tag1 :tag2
		pushArgs := append([]string{"origin"}, argsWithPrefix...)
		err := repo.Cmd("push", pushArgs...).Run()
		if err != nil {
			return err
		}

		colorp.Successf("Successfully deleted tags: %v\n", tagNames)
		return nil
	},
}
