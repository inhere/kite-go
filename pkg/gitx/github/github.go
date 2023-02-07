package github

import (
	"fmt"
)

const (
	ContentHost = "https://raw.githubusercontent.com"
	// ContentURLTpl eg: https://raw.githubusercontent.com/gookit/slog/master/README.md
	ContentURLTpl = "https://raw.githubusercontent.com/%s/%s/%s"
)

// GitHub config struct
type GitHub struct {
	// *gitx.Config

	// Token person token.
	Token string
}

// ContentURL for GitHub repo file
func ContentURL(repoPath, branch, filePath string) string {
	return fmt.Sprintf(ContentURLTpl, repoPath, branch, filePath)
}
