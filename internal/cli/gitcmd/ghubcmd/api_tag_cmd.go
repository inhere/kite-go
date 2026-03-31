package ghubcmd

import (
	"fmt"
	"strings"

	"github.com/gookit/color/colorp"
	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/show"
	"github.com/gookit/gitw/gitutil"
	"github.com/gookit/goutil/strutil"
	"github.com/inhere/kite-go/internal/app"
	"github.com/inhere/kite-go/internal/biz/cmdbiz"
	ghapi "github.com/inhere/kite-go/pkg/gitx/github"
)

// ApiCmd groups github api commands.
var ApiCmd = &gcli.Command{
	Name: "api",
	Desc: "github api helper commands",
	Subs: []*gcli.Command{
		ApiCommitCmd,
		ApiTagCmd,
	},
}

// ApiTagCmd groups github tag api commands.
var ApiTagCmd = &gcli.Command{
	Name: "tag",
	Desc: "github tag api commands",
	Subs: []*gcli.Command{
		ApiTagAddCmd,
		ApiTagListCmd,
	},
}

type apiTagAddOptions struct {
	cmdbiz.CommonOpts
	RepoPath    string
	SHA         string
	Message     string
	Version     string
	TaggerName  string
	TaggerEmail string
}

func (o *apiTagAddOptions) toCreateInput() (ghapi.TagCreateInput, error) {
	repoPath := strings.TrimSpace(o.RepoPath)
	if repoPath == "" {
		return ghapi.TagCreateInput{}, fmt.Errorf("the repository path is required, use --repo owner/repo")
	}

	_, _, err := gitutil.SplitPath(repoPath)
	if err != nil {
		return ghapi.TagCreateInput{}, err
	}

	version := strings.TrimSpace(o.Version)
	if version == "" {
		return ghapi.TagCreateInput{}, fmt.Errorf("the tag version is required")
	}

	if formatted, ok := gitutil.FormatVersion(version); ok {
		version = "v" + formatted
	}

	msg := strings.TrimSpace(o.Message)
	if msg == "" {
		return ghapi.TagCreateInput{}, fmt.Errorf("the tag message is required")
	}

	in := ghapi.TagCreateInput{
		RepoPath: repoPath,
		Version:  version,
		Message:  msg,
		Object:   strings.TrimSpace(o.SHA),
	}

	if o.TaggerName != "" || o.TaggerEmail != "" {
		in.Tagger = ghapi.Tagger{
			Name:  strings.TrimSpace(o.TaggerName),
			Email: strings.TrimSpace(o.TaggerEmail),
		}
	}
	return in, nil
}

func (o *apiTagAddOptions) resolveCreateInput(gh *ghapi.GitHub) (ghapi.TagCreateInput, error) {
	in, err := o.toCreateInput()
	if err != nil {
		return ghapi.TagCreateInput{}, err
	}

	if strings.TrimSpace(in.Object) != "" {
		return in, nil
	}

	info, err := gh.GetLatestCommit(in.RepoPath)
	if err != nil {
		return ghapi.TagCreateInput{}, err
	}

	in.Object = info.SHA
	return in, nil
}

var apiTagAddOpts = &apiTagAddOptions{}

// ApiTagAddOptionsForTestType exports test-only option shape.
type ApiTagAddOptionsForTestType = apiTagAddOptions

// ApiTagAddOptionsForTest provides test access for input normalization.
var ApiTagAddOptionsForTest ApiTagAddOptionsForTestType

// ToCreateInputForTest converts options to API input in tests.
func (o *ApiTagAddOptionsForTestType) ToCreateInputForTest() (ghapi.TagCreateInput, error) {
	return (*apiTagAddOptions)(o).toCreateInput()
}

// ApiTagAddCmd creates an annotated tag by GitHub API.
var ApiTagAddCmd = &gcli.Command{
	Name: "add",
	Desc: "create annotated tag by github api, similar to `git tag -a`",
	Help: `
# Examples:
  {$fullCmd} -r owner/repo --sha abc123 -v v1.2.3 -m "release v1.2.3"
  {$fullCmd} -r owner/repo --sha abc123 -v 1.2.3 -m "release v1.2.3" --tagger-name kite --tagger-email kite@example.com
`,
	Config: func(c *gcli.Command) {
		apiTagAddOpts.BindProxyConfirm(c)
		c.BoolOpt(&apiTagAddOpts.DryRun, "dry-run", "dry", false, "run workflow, but dont real execute command")
		c.StrOpt2(&apiTagAddOpts.RepoPath, "repo, r", "repository path, format: owner/repo")
		c.StrOpt(&apiTagAddOpts.SHA, "sha", "", "target commit sha")
		c.StrOpt2(&apiTagAddOpts.Version, "version, v", "tag version, eg: v2.0.1")
		c.StrOpt2(&apiTagAddOpts.Message, "message, m", "annotated tag message")
		c.StrOpt(&apiTagAddOpts.TaggerName, "tagger-name", "", "tagger name for annotated tag")
		c.StrOpt(&apiTagAddOpts.TaggerEmail, "tagger-email", "", "tagger email for annotated tag")
	},
	Func: func(c *gcli.Command, _ []string) error {
		gh := app.Ghub()
		if strutil.IsBlank(gh.Token) {
			return c.NewErr("github token is empty, please configure github.token or GITHUB_PA_TOKEN")
		}

		in, err := apiTagAddOpts.resolveCreateInput(gh)
		if err != nil {
			return err
		}

		show.AList("create github tag", map[string]any{
			"Repo":         in.RepoPath,
			"SHA":          in.Object,
			"Version":      in.Version,
			"Message":      in.Message,
			"Tagger Name":  in.Tagger.Name,
			"Tagger Email": in.Tagger.Email,
			"Dry Run":      apiTagAddOpts.DryRun,
		})

		if apiTagAddOpts.DryRun {
			return nil
		}

		resp, err := gh.CreateAnnotatedTag(in)
		if err != nil {
			return err
		}

		colorp.Successf("Successful create github tag: %s\n", in.Version)
		if resp.Ref != "" {
			colorp.Infof("Created ref: %s\n", resp.Ref)
		}
		if resp.Object.SHA != "" {
			colorp.Infof("Tag object SHA: %s\n", resp.Object.SHA)
		}
		return nil
	},
}
