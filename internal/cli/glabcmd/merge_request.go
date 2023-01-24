package glabcmd

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/show"
	"github.com/gookit/gitw/gitutil"
	"github.com/gookit/goutil/maputil"
	"github.com/gookit/goutil/strutil"
	"github.com/gookit/goutil/sysutil"
	"github.com/inhere/kite/app"
	"github.com/inhere/kite/pkg/gitx"
	"github.com/inhere/kite/pkg/gitx/gitlab"
)

var (
	mrOpts = struct {
		gitx.CommonOpts
		new bool
		// create pr link
		direct bool
		source string
		target string
		openBr string // same of target + openIt=true
		// open link on browser
		openIt bool
	}{}

	// MergeRequestCmd command
	MergeRequestCmd = &gcli.Command{
		Name:    "merge-request",
		Aliases: []string{"pr", "mr", "pull-request"},
		Desc:    "Create new merge requests(PR/MR) by given project information",
		Config: func(c *gcli.Command) {
			mrOpts.BindCommonFlags(c)

			c.BoolOpt(&mrOpts.new, "new", "", false,
				"Open new PR page link on browser. eg: http://my.gitlab.com/group/repo/merge_requests/new",
			)
			c.BoolOpt(&mrOpts.direct, "direct", "d", false,
				"The PR is direct from fork to main repository",
			)

			c.StrOpt(&mrOpts.openBr, "open", "o", "Set target branch and open PR link on browser")
			c.StrOpt(&mrOpts.source, "source", "s", "The source branch name, default is current `BRANCH`")
			c.StrOpt(&mrOpts.target, "target", "t", "The target branch name, default is current `BRANCH`")
			c.AddArg("repoPath", "The project name with path in self-host gitlab.\nif empty will fetch from workdir")
		},
		Help: `
Special:
  @, h, HEAD - Current branch name.
  @s         - Source branch name.
  @t         - Target branch name.
`,
		Examples: `
  {$binWithCmd}                       Will generate PR link for fork 'HEAD_BRANCH' to main 'HEAD_BRANCH'
  {$binWithCmd} -o h                  Will open PR link for fork 'HEAD_BRANCH' to main 'HEAD_BRANCH' on browser
  {$binWithCmd} -o qa                 Will open PR link for main 'HEAD_BRANCH' to main 'qa' on browser
  {$binWithCmd} -t qa                 Will generate PR link for main 'HEAD_BRANCH' to main 'qa'
  {$binWithCmd} -t qa --direct        Will generate PR link for fork 'HEAD_BRANCH' to main 'qa'
  # Will generate PR link for 'group/repo', from 'dev' to 'qa' branch
  {binWithCmd} -o dev -t qa group/repo
`,
		Func: func(c *gcli.Command, args []string) (err error) {
			workdir := c.WorkDir()
			repoPath := c.Arg("repoPath").String()

			var group, name string
			gl := app.Glab()
			glp := gitlab.NewGlProject(workdir, gl)

			if gl.HostUrl == "" {
				c.Println("TIP: gitlab.host_url is empty, try fetch from git remote")
				gl.HostUrl = glp.Repo().DefaultRemoteInfo().HTTPHost()
			}

			// http://my.gitlab.com/group/repo/merge_requests/new
			if mrOpts.new {
				if repoPath == "" {
					repoPath = glp.Repo().DefaultRemoteInfo().Path()
				}

				link := gl.HostUrl + "/" + repoPath
				return sysutil.OpenBrowser(link)
			}

			mrOpts.source, _ = glp.ResolveBranch(mrOpts.source)

			if mrOpts.openBr != "" {
				openBr, _ := glp.ResolveBranch(mrOpts.openBr)
				mrOpts.openIt = true
				mrOpts.target = openBr
			} else {
				mrOpts.target, _ = glp.ResolveBranch(mrOpts.target)
			}

			var srcPid, dstPid string
			var mrInfo *gitlab.PRLinkQuery
			if strutil.IsNotBlank(repoPath) {
				group, name, err = gitutil.SplitPath(repoPath)
				if err != nil {
					return err
				}

				srcPid = gitlab.BuildProjectID(group, name)
			} else {
				if err := glp.CheckForkRemote(); err != nil {
					return err
				}

				if mrOpts.target == "@s" {
					mrOpts.target = mrOpts.source
				}

				dstPid = glp.MainProjectId()
				// srcPid = glp.MainProjectId()
				if mrOpts.direct {
					srcPid = glp.ForkProjectId()
				}

				repoPath = glp.MainRmtInfo().Path()
			}

			show.AList("Some Options Info", maputil.Data{
				"direct from fork": mrOpts.direct,
				"open page link":   mrOpts.openIt,
			})

			mrInfo = gitlab.NewPRLinkQuery(srcPid, mrOpts.source, dstPid, mrOpts.target)
			mrInfo.RepoPath = repoPath

			show.AList("Merge Request Info", mrInfo)

			// link := glp.MargeRequestURL(mrInfo)
			link := mrInfo.BuildURL(gl.HostUrl)
			c.Infoln("Merge Request Link:")
			c.Println("   ", link)

			if mrOpts.openIt {
				err = sysutil.OpenBrowser(link)
			}
			return
		},
	}
)
