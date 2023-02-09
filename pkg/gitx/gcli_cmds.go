package gitx

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/gitw"
	"github.com/gookit/goutil/strutil"
	"github.com/gookit/goutil/sysutil"
)

// CommonOpts some common vars struct
type CommonOpts struct {
	Proxy   bool
	DryRun  bool
	Confirm bool
	Workdir string
	GitHost string
}

// BindCommonFlags for some git commands
func (co *CommonOpts) BindCommonFlags(c *gcli.Command) {
	c.BoolOpt(&co.DryRun, "dry-run", "dry", false, "run workflow, but dont real execute command")
	c.BoolOpt2(&co.Proxy, "proxy,p", "manual enable set proxy ENV config")
	c.StrOpt(&co.Workdir, "workdir", "w", "", "the command workdir path")
	c.BoolOpt2(&co.Confirm, "confirm", "confirm ask before executing command")
}

var orOpts = struct {
	remote string
}{}

// NewOpenRemoteCmd instance
func NewOpenRemoteCmd(hostUrlGetter func() string) *gcli.Command {
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
				var hostUrl string
				if hostUrlGetter != nil {
					hostUrl = hostUrlGetter()
				}

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
