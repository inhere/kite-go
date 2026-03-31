package github

import (
	"fmt"
	"net/url"
	"strings"
)

// CommitAuthor info from commit payload.
type CommitAuthor struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Date  string `json:"date"`
}

// CommitDetail info from commit payload.
type CommitDetail struct {
	Author  CommitAuthor `json:"author"`
	Message string       `json:"message"`
}

// CommitInfo represents a github commit item.
type CommitInfo struct {
	SHA     string       `json:"sha"`
	HTMLURL string       `json:"html_url"`
	Commit  CommitDetail `json:"commit"`
}

// GetLatestCommit gets latest commit info on the repository default branch.
func (g *GitHub) GetLatestCommit(repoPath string) (*CommitInfo, error) {
	repoPath = strings.TrimSpace(repoPath)
	if repoPath == "" {
		return nil, fmt.Errorf("github: repo path is required")
	}
	if strings.TrimSpace(g.Token) == "" {
		return nil, fmt.Errorf("github: token is required")
	}

	var items []CommitInfo
	err := g.getJSON("/repos/"+strings.Trim(repoPath, "/")+"/commits", url.Values{
		"per_page": []string{"1"},
	}, &items)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("github: no commits found for repo %s", repoPath)
	}
	return &items[0], nil
}
