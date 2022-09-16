package gitlab

import "github.com/gookit/gcli/v3"

var (
	mrOpts = struct {
		new    bool
		direct bool
		open   string
		source string
		target string
	}{}

	MergeRequest = &gcli.Command{
		Name:    "merge-request",
		Aliases: []string{"pr", "mr", "pull-request"},
		Desc:    "Generate an PR/MR link for given project information",
		Config: func(c *gcli.Command) {
			bindCommonFlags(c)

			c.BoolOpt(&mrOpts.new, "new", "", false,
				"Open new pr page on browser. eg: http://my.gitlab.com/group/repo/merge_requests/new",
			)

			c.BoolOpt(&mrOpts.direct, "direct", "d", false,
				"The PR is direct from fork to main repository",
			)
			c.StrOpt(&mrOpts.open, "open", "o", "", "generate PR to `BRANCH` and open link on browser")
			c.StrOpt(&mrOpts.source, "source", "s", "", "The source branch name, default is current `BRANCH`")
			c.StrOpt(&mrOpts.target, "target", "t", "", "The target branch name, default is current `BRANCH`")
		},
		Help: `
Special:
  @, HEAD - Current branch.
  @s      - Source branch.
  @t      - Target branch.
`,
		Examples: `
  {$binWithCmd}                       Will generate PR link for fork 'HEAD_BRANCH' to main 'HEAD_BRANCH'
  {$binWithCmd} -o @                  Will open PR link for fork 'HEAD_BRANCH' to main 'HEAD_BRANCH' on browser
  {$binWithCmd} -o qa                 Will open PR link for main 'HEAD_BRANCH' to main 'qa' on browser
  {$binWithCmd} -t qa                 Will generate PR link for main 'HEAD_BRANCH' to main 'qa'
  {$binWithCmd} -t qa --direct       Will generate PR link for fork 'HEAD_BRANCH' to main 'qa'
`,
		Func: func(c *gcli.Command, args []string) error {

			return nil
		},
	}
)
