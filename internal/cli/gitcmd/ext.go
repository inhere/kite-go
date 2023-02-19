package gitcmd

import (
	"fmt"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/gflag"
	"github.com/gookit/gitw"
	"github.com/gookit/gitw/gitutil"
	"github.com/gookit/gitw/gmoji"
	"github.com/gookit/validate"
	"github.com/inhere/kite/internal/apputil"
	"github.com/inhere/kite/internal/biz/cmdbiz"
	"github.com/inhere/kite/pkg/gitx"
)

// NewCloneCmd instance
func NewCloneCmd(cfgGetter gitx.ConfigProviderFn) *gcli.Command {
	var clOpts = struct {
		cmdbiz.CommonOpts
		gitProto bool
		// arg
		repoPath  string
		localName string
	}{}

	return &gcli.Command{
		Name:    "clone",
		Aliases: []string{"down"},
		Desc:    "quick lone an git repository to local",
		Config: func(c *gcli.Command) {
			clOpts.BindCommonFlags(c)

			c.BoolOpt2(&clOpts.gitProto, "git, g", "Use git protocol for git clone")

			c.AddArg("repoPath", "repo path or full remote url", true).WithAfterFn(func(a *gflag.CliArg) error {
				clOpts.repoPath = a.String()
				return nil
			})
			c.AddArg("local", "custom the local name of clone repository").WithAfterFn(func(a *gflag.CliArg) error {
				clOpts.localName = a.String()
				return nil
			})
		},
		Func: func(c *gcli.Command, args []string) error {
			var remoteURL string
			if validate.IsFullURL(clOpts.repoPath) {
				remoteURL = clOpts.repoPath
			} else if gitutil.IsRepoPath(clOpts.repoPath) {
				cfg := cfgGetter()
				remoteURL = cfg.HostUrl + "/" + clOpts.repoPath
			} else if ghURL, ok := gitutil.ResolveGhURL(clOpts.repoPath); ok {
				remoteURL = ghURL
			}

			if remoteURL != "" {
				return gitw.
					New("clone", clOpts.repoPath).
					ArgIf(clOpts.localName, clOpts.localName != "").
					PrintCmdline().
					Run()
			}

			return c.NewErrf("cannot clone repository from the %s", clOpts.repoPath)
		},
	}
}

// NewGitEmojisCmd instance
func NewGitEmojisCmd() *gcli.Command {
	var geOpts = struct {
		list   bool
		lang   string
		render string
		output string
		search gcli.String
	}{}

	return &gcli.Command{
		Name:    "emoji",
		Desc:    "checkout an new branch for development from `dist` remote",
		Aliases: []string{"moji"},
		Config: func(c *gcli.Command) {
			c.BoolOpt2(&geOpts.list, "list, ls", "list all git emojis")
			c.StrOpt2(&geOpts.lang, "lang,l", "the language for git emojis, allow: en, zh-CN", gflag.WithDefault(gmoji.LangEN))
			c.VarOpt(&geOpts.search, "search", "s", "search git emojis by `keywords`, multi by comma ','")
			c.StrOpt2(&geOpts.render, "render, r", "render emoji code to emojis of input message")
			c.StrOpt2(&geOpts.output, "output, o", "the output after rendered, default is stdout")
		},
		Func: func(c *gcli.Command, args []string) error {
			em, err := gmoji.Emojis(geOpts.lang)
			if err != nil {
				return err
			}

			if geOpts.list {
				c.Warnf("All git emojis(total: %d):\n", em.Len())
				fmt.Println(em.String())
				return nil
			}

			if geOpts.search != "" {
				sub := em.Search(geOpts.search.Strings(), 10)

				c.Warnf("Matched emojis(total: %d):\n", sub.Len())
				fmt.Println(sub.String())
				return nil
			}

			if geOpts.render != "" {
				src := apputil.ReadSource(geOpts.render)
				src = em.RenderCodes(src)

				fmt.Println(src)
				return nil
			}

			return c.NewErr("please input an option for operation.")
		},
	}
}
