package gitcmd

import (
	"fmt"
	"strings"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/show"
	"github.com/gookit/gitw/gitutil"
	"github.com/gookit/goutil/maputil"
	"github.com/gookit/goutil/strutil"
	"github.com/gookit/goutil/sysutil"
	"github.com/inhere/kite-go/internal/app"
	"github.com/inhere/kite-go/internal/biz/cmdbiz"
)

var (
	mrOpts = struct {
		cmdbiz.CommonOpts
		new bool
		// create pr link
		direct bool
		source string // source branch
		target string // target branch
		openBr string // same of target + openIt=true
		// open link on browser
		openIt bool
		// auto search .git repo on parent dir
		// findRepo bool
	}{}
)

// NewPullRequestCmd command
func NewPullRequestCmd() *gcli.Command {
	pr := &gcli.Command{
		Name:    "pull-request",
		Aliases: []string{"pr", "mr", "merge-request"},
		Desc:    "Create new merge requests(PR/MR) by given project information",
		Config: func(c *gcli.Command) {
			mrOpts.BindCommonFlags(c)

			c.BoolOpt(&mrOpts.new, "new", "", false,
				"Open new PR page link on browser. eg: http://my.git.com/group/repo/merge_requests/new",
			)
			c.BoolOpt2(&mrOpts.direct, "direct,d", "The PR is direct from fork to main repository")
			// c.BoolOpt2(&mrOpts.findRepo, "find-repo, find", "auto find repo .git dir on parent dirs")

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
  {$binWithCmd} -o dev -t qa group/repo
`,
		Func: mergeRequestHandle,
	}

	return pr
}

func mergeRequestHandle(c *gcli.Command, _ []string) (err error) {
	repoDir := mrOpts.Workdir

	gx := app.Gitx()
	gp := gx.LoadRepo(repoDir)

	hostUrl := gx.HostUrl
	if gp.IsGitRepo() {
		c.Infoln("TIP: in git repository, try fetch host_url from default remote")
		hostUrl = gp.DefaultRemoteInfo().HTTPHost(gx.DisableHTTPS)
	}

	// http://my.gitlab.com/group/repo/merge_requests/new
	repoPath := c.Arg("repoPath").String()
	if mrOpts.new {
		if repoPath == "" {
			repoPath = gp.DefaultRemoteInfo().Path()
		}

		link := hostUrl + "/" + repoPath + "/pulls"
		return sysutil.OpenBrowser(link)
	}

	mrOpts.source, _ = gp.ResolveBranch(mrOpts.source)
	if mrOpts.openBr != "" {
		openBr, _ := gp.ResolveBranch(mrOpts.openBr)
		mrOpts.openIt = true
		mrOpts.target = openBr
	} else {
		mrOpts.target, _ = gp.ResolveBranch(mrOpts.target)
	}

	var srcPid, dstPid string
	// var group, name string
	if strutil.IsNotBlank(repoPath) {
		repoPath = strings.TrimSpace(repoPath)
		_, _, err = gitutil.SplitPath(repoPath)
		if err != nil {
			return err
		}

		srcPid, dstPid = repoPath, repoPath
	} else {
		if err := gp.CheckRemote(); err != nil {
			return err
		}

		if mrOpts.target == "@s" {
			mrOpts.target = mrOpts.source
		}

		dstPid = gp.SrcRemoteInfo().Path()
		srcPid = gp.DefaultRemoteInfo().Path()
		if mrOpts.direct {
			dstPid = gp.DefaultRemoteInfo().Path()
		}

		if mrOpts.direct || mrOpts.target == mrOpts.source {
			repoPath = gp.DefaultRemoteInfo().Path()
		} else {
			repoPath = dstPid
		}
	}

	show.AList("Current Information", maputil.Data{
		"Direct from fork":  mrOpts.direct,
		"Open browser link": mrOpts.openIt,
		"From group/name":   srcPid,
		"Into group/name":   dstPid,
		"Into branch name":  mrOpts.target,
	})

	// http://git.your.com/GROUP/NAME/compare/TARGET_BRANCH...GROUP1/NAME:SOURCE_BRANCH
	link := fmt.Sprintf(
		"%s/%s/compare/%s...%s:%s",
		hostUrl, dstPid, mrOpts.target, srcPid, mrOpts.source,
	)
	c.Warnln("Pull Request Link:")
	c.Println("  ", link)

	if mrOpts.openIt {
		err = sysutil.OpenBrowser(link)
	}
	return
}
