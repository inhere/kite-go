package github

import (
	"fmt"

	"github.com/gookit/gitw"
	"github.com/inhere/kite-go/pkg/gitx"
)

const (
	ContentHost = "https://raw.githubusercontent.com"
	// ContentURLTpl eg: https://raw.githubusercontent.com/gookit/slog/master/README.md
	ContentURLTpl = "https://raw.githubusercontent.com/%s/%s/%s"
)

// GitHub config struct
type GitHub struct {
	*gitx.Config

	// GitHub 文件, Releases, archive, gist, raw.githubusercontent.com 文件代理加速下载服务.
	// eg: https://ghproxy.com
	ProxyHost string `json:"proxy_host"`
	// Username on https://github.com
	Username string `json:"username"`
	// Token person token.
	Token string `json:"token"`
	// BaseApi api url
	BaseApi string `json:"base_api"`
}

// New config instance
func New(cfg *gitx.Config) *GitHub {
	cfg.HostType = gitx.HostGitHub

	if cfg.HostUrl == "" {
		cfg.HostUrl = gitw.GitHubURL
	}
	if cfg.GitUrl == "" {
		cfg.GitUrl = gitw.GitHubGit
	}

	return &GitHub{Config: cfg}
}

// ContentURL for GitHub repo file
func ContentURL(repoPath, branch, filePath string) string {
	return fmt.Sprintf(ContentURLTpl, repoPath, branch, filePath)
}
