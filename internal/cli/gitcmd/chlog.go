package gitcmd

import (
	"path/filepath"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/gitw"
	"github.com/gookit/gitw/chlog"
	"github.com/gookit/goutil/dump"
)

var (
	chlogOpts = struct {
		limit     gcli.String
		sha1      string
		sha2      string
		style     string
		repoUrl   string
		dstFile   string
		noMerges  bool
		unShallow bool
		fetchTags bool
	}{}

	ChangelogCmd = &gcli.Command{
		Name:    "chlog",
		Desc:    "batch pull multi git directory by `git pull`",
		Aliases: []string{"cl", "clog", "changelog"},
		Examples: `
  {$binWithCmd} last head
  {$binWithCmd} last head --style gh-release --no-merges
  {$binWithPath} v2.0.9 v2.0.10 --no-merges --style gh-release --exclude "cs-fixer,format codes"
`,
		Config: func(c *gcli.Command) {
			c.AddArg("oldVersion", `The old version. eg: v1.0.2, 349238b
- keywords 'last/latest' will auto use latest tag
- keywords 'prev/previous' will auto use previous tag`).
				WithFn(func(arg *gcli.Argument) {
					arg.Required = true
				})

			c.AddArg("newVersion", `The new version. eg: v1.2.2, 66c0df1
- keywords 'head' will auto use Head commit`).
				WithFn(func(arg *gcli.Argument) {
					arg.Required = true
				})

			c.VarOpt(&chlogOpts.limit, "limit", "", "limit update the given dir names")
			c.StrOpt2(&chlogOpts.dstFile, "file", "Export changelog message to the file, default dump to stdout")
			c.StrOpt2(&chlogOpts.repoUrl, "repo-url", `
The git repo URL address. eg: https://github.com/inhere/kite
 default will auto use current git origin remote url
`)
			c.StrOpt(&chlogOpts.style, "style", "s", "default", `
The style for generate for changelog.
 allow: markdown(md), simple, gh-release(ghr)
`)
			c.BoolOpt2(&chlogOpts.fetchTags, "fetch-tags", "Update repo tags list by 'git fetch --tags'")
			c.BoolOpt2(&chlogOpts.noMerges, "no-merges", "dont contains merge request logs")
			c.BoolOpt2(&chlogOpts.unShallow, "unshallow", "Convert to a complete warehouse, useful on GitHub Action.")
		},
		Func: func(c *gcli.Command, args []string) error {
			baseDir := c.Arg("baseDir").String()
			absDir, err := filepath.Abs(baseDir)
			if err != nil {
				return err
			}

			repo := gitw.NewRepo(absDir)
			dump.P(repo.DefaultRemoteInfo())

			cl := chlog.New()
			cl.Formatter = &chlog.MarkdownFormatter{
				RepoURL: "https://github.com/gookit/gitw",
			}
			// with some settings ...
			cl.WithConfigFn(func(c *chlog.Config) {
				c.GroupPrefix = "\n### "
				c.GroupSuffix = "\n"
			})

			chlogOpts.sha1 = c.Arg("oldVersion").String()
			chlogOpts.sha2 = c.Arg("newVersion").String()
			dump.P(chlogOpts)

			// fetch git log
			cl.FetchGitLog(chlogOpts.sha1, chlogOpts.sha2, "--no-merges")

			if err = cl.Generate(); err != nil {
				return err
			}

			dump.P(cl.Changelog())
			return nil
		},
	}
)
