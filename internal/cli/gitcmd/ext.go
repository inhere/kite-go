package gitcmd

import (
	"fmt"
	"strings"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/gflag"
	"github.com/gookit/gitw"
	"github.com/gookit/gitw/gitutil"
	"github.com/gookit/gitw/gmoji"
	"github.com/gookit/goutil/strutil"
	"github.com/gookit/goutil/sysutil"
	"github.com/inhere/kite/internal/apputil"
	"github.com/inhere/kite/internal/biz/cmdbiz"
	"github.com/inhere/kite/pkg/gitx"
)

// NewCloneCmd instance
func NewCloneCmd(cfgGetter gitx.ConfigProviderFn) *gcli.Command {
	var clOpts = struct {
		cmdbiz.CommonOpts
		useSsh bool
		// arg
		repoPath  string
		localName string
	}{}

	return &gcli.Command{
		Name: "clone",
		// Aliases: []string{"down"},
		Desc: "quick lone an git repository to local",
		Config: func(c *gcli.Command) {
			clOpts.BindCommonFlags(c)

			c.BoolOpt2(&clOpts.useSsh, "git, ssh, g", "Use ssh protocol for git clone")
			c.AddArg("repoPath", "repo path or full remote repository url", true).WithAfterFn(func(a *gflag.CliArg) error {
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
			repoPath := clOpts.repoPath

			if gitutil.IsFullURL(repoPath) {
				remoteURL = repoPath
			} else if gitutil.IsRepoPath(repoPath) {
				cfg := cfgGetter()
				confMp := apputil.CmdConfigData(cfg.HostType, c.Name)
				if confMp.Bool("use_ssh") {
					clOpts.useSsh = true
				}

				c.Infof("TIP: only input a repoPath %q, will clone from: %s\n", repoPath, cfg.HostUrl)
				remoteURL = cfg.BuildRepoURL(repoPath, clOpts.useSsh)
			} else if ghURL, ok := gitutil.ResolveGhURL(repoPath); ok {
				c.Infof("TIP: input repoPath is start withs %s, will clone from GitHub\n", gitw.GitHubHost)
				remoteURL = ghURL
			}

			if remoteURL != "" {
				return gitw.
					New("clone", remoteURL).
					ArgIf(clOpts.localName, clOpts.localName != "").
					WithWorkDir(clOpts.Workdir).
					WithDryRun(clOpts.DryRun).
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
		search gflag.String
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
				src, err := apputil.ReadSource(geOpts.render)
				if err != nil {
					return err
				}

				fmt.Println(em.RenderCodes(src))
				return nil
			}

			return c.NewErr("please input an option for operation.")
		},
	}
}

var orOpts = struct {
	source   bool
	remote   string
	repoPath string
}{}

// NewOpenRemoteCmd instance
func NewOpenRemoteCmd(cfgGetter gitx.ConfigProviderFn) *gcli.Command {
	return &gcli.Command{
		Name: "open",
		Desc: "open the git remote repo address on browser",
		Config: func(c *gcli.Command) {
			c.StrOpt(&orOpts.remote, "remote", "r", "the remote name, if not input will use default remote")
			c.BoolOpt2(&orOpts.source, "source,src, s", "direct open the source_remote repository")

			c.AddArg("repoPath", "the git repo path with name. format: GROUP/NAME").WithAfterFn(func(a *gflag.CliArg) error {
				orOpts.repoPath = a.String()
				return nil
			})
		},
		Func: func(c *gcli.Command, args []string) error {
			cfg := cfgGetter()

			remote := orOpts.remote
			repoPath := orOpts.repoPath
			hostUrl := cfg.HostUrl

			var repoUrl string
			if orOpts.source {
				remote = cfg.SourceRemote
			}

			if strutil.IsNotBlank(repoPath) {
				// special github url
				if strings.HasPrefix(repoPath, gitx.GitHubHost) {
					repoUrl = "https://" + repoPath
				} else if hostUrl != "" {
					repoUrl = hostUrl + "/" + repoPath
				} else {
					repo := gitw.NewRepo(c.WorkDir())
					repoUrl = repo.RemoteInfo(remote).HTTPHost() + "/" + repoPath
				}
			} else {
				// parse from git repo
				repo := gitw.NewRepo(c.WorkDir())
				rmt := repo.RemoteInfo(remote)

				if hostUrl != "" {
					repoUrl = hostUrl + "/" + rmt.RepoPath()
				} else {
					repoUrl = rmt.URLOrBuild()
				}
			}

			c.Infoln("Open URL:", repoUrl)
			return sysutil.OpenBrowser(repoUrl)
		},
	}
}
