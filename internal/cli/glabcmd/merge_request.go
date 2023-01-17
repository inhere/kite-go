package glabcmd

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/gitw"
	"github.com/gookit/gitw/gitutil"
	"github.com/gookit/goutil/dump"
	"github.com/gookit/goutil/strutil"
	"github.com/inhere/kite/pkg/gitx"
)

var (
	mrOpts = struct {
		gitx.CommonOpts
		new    bool
		direct bool
		open   string
		source string
		target string
	}{}

	// MergeRequestCmd command
	MergeRequestCmd = &gcli.Command{
		Name:    "merge-request",
		Aliases: []string{"pr", "mr", "pull-request"},
		Desc:    "Create new merge requests(PR/MR) by given project information",
		Config: func(c *gcli.Command) {
			mrOpts.BindCommonFlags(c)

			c.BoolOpt(&mrOpts.new, "new", "", false,
				"Open new pr page on browser. eg: http://my.gitlab.com/group/repo/merge_requests/new",
			)
			c.BoolOpt(&mrOpts.direct, "direct", "d", false,
				"The PR is direct from fork to main repository",
			)

			c.StrOpt(&mrOpts.source, "source", "s", "The source branch name, default is current `BRANCH`")
			c.StrOpt(&mrOpts.open, "open", "o", "generate PR to `BRANCH` and open link on browser")
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

			if strutil.IsNotBlank(repoPath) {
				group, name, err = gitutil.SplitPath(repoPath)
				if err != nil {
					return err
				}
			} else {
				repo := gitw.NewRepo(workdir)

				rtInfo := repo.DefaultRemoteInfo()
				group, name = rtInfo.Group, rtInfo.Repo
			}

			dump.P(group, name)

			return nil
		},
	}
)
