package gitx

import (
	"github.com/gookit/gitw"
	"github.com/gookit/goutil/maputil"
)

const (
	HostGitHub  = "github"
	HostGitlab  = "gitlab"
	HostDefault = "git"
)

// ConfigProviderFn type
type ConfigProviderFn func() *Config

// Config for gitx
type Config struct {
	// GitUrl ssh git host url
	GitUrl string `json:"git_url"`
	// HostUrl http host url. eg: https://gitlab.myself.com
	HostUrl string `json:"host_url"`
	// HostType eg: HostGitHub
	HostType string `json:"host_type"`
	// AutoSign auto add author sign on commit
	AutoSign bool `json:"auto_sign"`
	// ForkMode enable git fork mode for develop.
	// If is False, use branch mode, will ignore SourceRemote setting.
	ForkMode bool `json:"fork_mode"`
	// DisableHTTPS disable https(eg: on PR/MR link). if is true, will use HTTP url.
	DisableHTTPS bool `json:"disable_https"`
	// SourceRemote the source remote name, it is center repo.
	SourceRemote string `json:"source_remote"`
	// DefaultRemote the default upstream remote name, use for develop. default: origin.
	// It should be forked from SourceRemote.
	DefaultRemote string `json:"default_remote"`
	// DefaultBranch name, default is gitw.DefaultBranchName
	DefaultBranch string `json:"default_branch"`
	// BranchAliases branch aliases
	BranchAliases maputil.Aliases `json:"branch_aliases"`
	// PrUrlFormat pull request URL format template. can use var like {host}
	PrUrlFormat string `json:"pr_url_format"`
}

// NewConfig instance
func NewConfig() *Config {
	return &Config{
		HostType: HostDefault,
		// remote
		SourceRemote:  DefaultSrcRemote,
		DefaultRemote: gitw.DefaultRemoteName,
		// branch
		DefaultBranch: gitw.DefaultBranchName,
	}
}

// LoadRepo by given git repo dir.
func (c *Config) LoadRepo(repoDir string) *GitLoc {
	return NewGitLoc(repoDir, c)
}

// Init config.
func (c *Config) Init() *Config {
	if c.DefaultBranch == "" {
		c.DefaultBranch = gitw.DefaultBranchName
	}
	if c.DefaultRemote == "" {
		c.DefaultRemote = gitw.DefaultRemoteName
	}

	return c
}

func (c *Config) IsDefaultRemote(remote string) bool {
	return c.DefaultRemote == remote
}

func (c *Config) IsSourceRemote(remote string) bool {
	return c.IsCenterRemote(remote)
}

func (c *Config) IsCenterRemote(remote string) bool {
	return c.SourceRemote == remote
}

// IsForkMode check
func (c *Config) IsForkMode() bool {
	return c.ForkMode && c.SourceRemote != ""
}

// OriginBranch name build
func (c *Config) OriginBranch(br string) string {
	return c.DefaultRemote + "/" + br
}

// SourceBranch name build
func (c *Config) SourceBranch(br string) string {
	return c.SourceRemote + "/" + br
}

// BuildRepoURL build
func (c *Config) BuildRepoURL(repoPath string, useSsh bool) string {
	if useSsh {
		return c.GitUrl + ":" + repoPath
	}
	return c.HostUrl + "/" + repoPath
}

// ResolveBranch branch name
func (c *Config) ResolveBranch(name string) string {
	return c.BranchAliases.ResolveAlias(name)
}

// ResolveAlias branch name
func (c *Config) ResolveAlias(name string) string {
	return c.ResolveBranch(name)
}

// Clone new config instance
func (c *Config) Clone() *Config {
	c1 := *c
	return &c1
}
