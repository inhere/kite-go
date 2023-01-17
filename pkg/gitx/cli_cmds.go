package gitx

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/gitw"
	"github.com/gookit/goutil/strutil"
	"github.com/gookit/goutil/sysutil"
)

// CommonOpts some common vars struct
type CommonOpts struct {
	DryRun  bool
	Workdir string
	GitHost string
}

// BindCommonFlags for some git commands
func (co CommonOpts) BindCommonFlags(c *gcli.Command) {
	c.BoolOpt(&co.DryRun, "dry-run", "", false, "run workflow, but dont real execute command")
	c.StrOpt(&co.Workdir, "workdir", "w", "", "the command workdir path")
}

var orOpts = struct {
	remote string
}{}

// NewOpenRemoteCmd instance
func NewOpenRemoteCmd(hostUrl string) *gcli.Command {
	return &gcli.Command{
		Name: "open",
		Desc: "open the git remote repo address on browser",
		Config: func(c *gcli.Command) {
			c.StrOpt(&orOpts.remote, "remote", "r", "the remote name, if not input will use default remote")
			c.AddArg("repoPath", "the git repo path with name. format: GROUP/NAME")
		},
		Func: func(c *gcli.Command, args []string) error {
			remote := orOpts.remote
			repoPath := c.Arg("repoPath").String()

			var repoUrl string
			if strutil.IsNotBlank(repoPath) {
				if hostUrl != "" {
					repoUrl = hostUrl + "/" + repoPath
				} else {
					repo := gitw.NewRepo(c.WorkDir())
					repoUrl = repo.RemoteInfo(remote).HTTPHost() + "/" + repoPath
				}
			} else {
				repo := gitw.NewRepo(c.WorkDir())
				repoUrl = repo.RemoteInfo(remote).URLOrBuild()
			}

			c.Infoln("Open URL:", repoUrl)
			return sysutil.OpenBrowser(repoUrl)
		},
	}

}

// OpenRemoteRepo address
var OpenRemoteRepo = &gcli.Command{
	Name: "open",
	Desc: "open the git remote repo address on browser",
	Config: func(c *gcli.Command) {
		c.StrOpt(&orOpts.remote, "remote", "r", "the remote name, if not input will use default remote")
		c.AddArg("repoPath", "the git repo path with name. format: GROUP/NAME")
	},
	Func: func(c *gcli.Command, args []string) error {
		repo := gitw.NewRepo(c.WorkDir())

		url := repo.DefaultRemoteInfo().URLOrBuild()
		c.Infoln("Open URL:", url)

		return sysutil.OpenBrowser(url)
	},
}
