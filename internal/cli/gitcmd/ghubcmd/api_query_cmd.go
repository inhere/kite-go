package ghubcmd

import (
	"fmt"
	"strings"

	"github.com/gookit/cliui/show"
	"github.com/gookit/color/colorp"
	"github.com/gookit/gcli/v3"
	"github.com/gookit/gitw/gitutil"
	"github.com/gookit/goutil/strutil"
	"github.com/inhere/kite-go/internal/app"
	ghapi "github.com/inhere/kite-go/pkg/gitx/github"
)

// ApiCommitCmd groups github commit api commands.
var ApiCommitCmd = &gcli.Command{
	Name: "commit",
	Desc: "github commit api commands",
	Subs: []*gcli.Command{
		ApiCommitLatestCmd,
	},
}

type apiLatestCommitOptions struct {
	RepoPath string
}

func (o apiLatestCommitOptions) toInput() (string, error) {
	repoPath := strings.TrimSpace(o.RepoPath)
	if repoPath == "" {
		return "", fmt.Errorf("the repository path is required, use --repo owner/repo")
	}

	_, _, err := gitutil.SplitPath(repoPath)
	if err != nil {
		return "", err
	}
	return repoPath, nil
}

var apiLatestCommitOpts = &apiLatestCommitOptions{}

// ApiCommitLatestCmd gets latest commit info by github api.
var ApiCommitLatestCmd = &gcli.Command{
	Name: "latest",
	Desc: "get latest commit info by github api",
	Help: `
# Examples:
  {$fullCmd} -r owner/repo
`,
	Config: func(c *gcli.Command) {
		c.StrOpt2(&apiLatestCommitOpts.RepoPath, "repo, r", "repository path, format: owner/repo")
	},
	Func: func(c *gcli.Command, _ []string) error {
		repoPath, err := apiLatestCommitOpts.toInput()
		if err != nil {
			return err
		}

		gh := app.Ghub()
		if strutil.IsBlank(gh.Token) {
			return c.NewErr("github token is empty, please configure github.token or GITHUB_PA_TOKEN")
		}

		info, err := gh.GetLatestCommit(repoPath)
		if err != nil {
			return err
		}

		show.AList("latest github commit", map[string]any{
			"Repo":    repoPath,
			"SHA":     info.SHA,
			"Author":  info.Commit.Author.Name,
			"Email":   info.Commit.Author.Email,
			"Date":    info.Commit.Author.Date,
			"Message": info.Commit.Message,
			"URL":     info.HTMLURL,
		})
		return nil
	},
}

type apiTagListOptions struct {
	RepoPath string
	Limit    int
}

func (o apiTagListOptions) toInput() (ghapi.TagListInput, error) {
	repoPath := strings.TrimSpace(o.RepoPath)
	if repoPath == "" {
		return ghapi.TagListInput{}, fmt.Errorf("the repository path is required, use --repo owner/repo")
	}

	_, _, err := gitutil.SplitPath(repoPath)
	if err != nil {
		return ghapi.TagListInput{}, err
	}

	limit := o.Limit
	switch {
	case limit <= 0:
		limit = 20
	case limit > 100:
		limit = 100
	}

	return ghapi.TagListInput{
		RepoPath: repoPath,
		Limit:    limit,
	}, nil
}

var apiTagListOpts = &apiTagListOptions{}

// ApiTagListCmd lists tags by github api.
var ApiTagListCmd = &gcli.Command{
	Name:    "list",
	Aliases: []string{"ls"},
	Desc:    "list tags by github api",
	Help: `
# Examples:
  {$fullCmd} -r owner/repo
  {$fullCmd} -r owner/repo --limit 10
`,
	Config: func(c *gcli.Command) {
		c.StrOpt2(&apiTagListOpts.RepoPath, "repo, r", "repository path, format: owner/repo")
		c.IntOpt2(&apiTagListOpts.Limit, "limit, l", "max tags count for listing")
	},
	Func: func(c *gcli.Command, _ []string) error {
		in, err := apiTagListOpts.toInput()
		if err != nil {
			return err
		}

		gh := app.Ghub()
		if strutil.IsBlank(gh.Token) {
			return c.NewErr("github token is empty, please configure github.token or GITHUB_PA_TOKEN")
		}

		items, err := gh.ListTags(in)
		if err != nil {
			return err
		}

		if len(items) == 0 {
			colorp.Infof("No tags found for the repository: %s\n", in.RepoPath)
			return nil
		}

		colorp.Infof("Tags for the repository %s:\n", in.RepoPath)
		for _, item := range items {
			line := "  " + item.Name
			if item.Commit.SHA != "" {
				line += "  " + item.Commit.SHA
			}
			colorp.Cyanf("%s\n", line)
		}
		return nil
	},
}
