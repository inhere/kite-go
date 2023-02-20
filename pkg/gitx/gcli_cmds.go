package gitx

import (
	"strings"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/gflag"
	"github.com/gookit/gitw"
	"github.com/gookit/goutil/strutil"
	"github.com/gookit/goutil/sysutil"
)

var orOpts = struct {
	remote   string
	repoPath string
}{}

// NewOpenRemoteCmd instance
func NewOpenRemoteCmd(hostUrlGetter func() string) *gcli.Command {
	return &gcli.Command{
		Name: "open",
		Desc: "open the git remote repo address on browser",
		Config: func(c *gcli.Command) {
			c.StrOpt(&orOpts.remote, "remote", "r", "the remote name, if not input will use default remote")
			c.AddArg("repoPath", "the git repo path with name. format: GROUP/NAME").WithAfterFn(func(a *gflag.CliArg) error {
				orOpts.repoPath = a.String()
				return nil
			})
		},
		Func: func(c *gcli.Command, args []string) error {
			remote := orOpts.remote
			// repoPath := c.Arg("repoPath").String()
			repoPath := orOpts.repoPath

			var hostUrl, repoUrl string
			if hostUrlGetter != nil {
				hostUrl = hostUrlGetter()
			}

			if strutil.IsNotBlank(repoPath) {
				// special github url
				if strings.HasPrefix(repoPath, GitHubHost) {
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
