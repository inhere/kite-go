package gitlab

import (
	"net/url"
	"strings"

	"github.com/gookit/gitw"
	"github.com/gookit/goutil/errorx"
)

// GlProject struct
type GlProject struct {
	*GitLab
	err error
	// dir for the project
	dir  string
	repo *gitw.Repo
	// remote info
	mainRemoteInfo *gitw.RemoteInfo
	forkRemoteInfo *gitw.RemoteInfo

	// other
	mainProjId, forkProjId string
}

// NewGlProject instance
func NewGlProject(dir string, gl *GitLab) *GlProject {
	return &GlProject{
		dir:    dir,
		GitLab: gl,
	}
}

// Err get
func (p *GlProject) Err() error {
	return p.err
}

func (p *GlProject) Repo() *gitw.Repo {
	if p.repo == nil {
		p.repo = gitw.NewRepo(p.dir)
	}
	return p.repo
}

func (p *GlProject) CheckRemote() error {
	if err := p.CheckForkRemote(); err != nil {
		return err
	}
	return p.CheckMainRemote()
}

func (p *GlProject) CheckForkRemote() error {
	defRmt := p.Repo().DefaultRemoteInfo()

	if !p.Repo().HasRemote(p.DefaultRemote) {
		return errorx.Newf("the fork remote '%s' is not found on %s", p.DefaultRemote, defRmt.Path())
	}
	return nil
}

func (p *GlProject) CheckMainRemote() error {
	defRmt := p.Repo().DefaultRemoteInfo()

	if !p.Repo().HasRemote(p.UpstreamRemote) {
		return errorx.Newf("the main remote '%s' is not found on %s", p.UpstreamRemote, defRmt.Path())
	}
	return nil
}

func (p *GlProject) ForkRmtInfo() *gitw.RemoteInfo {
	if p.forkRemoteInfo == nil {
		p.addError(p.CheckForkRemote())
		p.forkRemoteInfo = p.Repo().RemoteInfo(p.DefaultRemote)
	}
	return p.forkRemoteInfo
}

func (p *GlProject) MainRmtInfo() *gitw.RemoteInfo {
	if p.mainRemoteInfo == nil {
		p.addError(p.CheckForkRemote())
		p.mainRemoteInfo = p.Repo().RemoteInfo(p.UpstreamRemote)
	}
	return p.mainRemoteInfo
}

// ResolveBranch name
func (p *GlProject) ResolveBranch(brName string) (string, bool) {
	switch strings.ToUpper(brName) {
	case "", "@", "H", "HEAD":
		return p.Repo().CurBranchName(), true
	}
	return p.ResolveAlias(brName), false
}

func (p *GlProject) MainProjectId() string {
	if p.mainProjId == "" {
		p.mainProjId = url.PathEscape(p.MainRmtInfo().Path())
	}
	return p.mainProjId
}

func (p *GlProject) ForkProjectId() string {
	if p.forkProjId == "" {
		p.forkProjId = url.PathEscape(p.ForkRmtInfo().Path())
	}
	return p.forkProjId
}

func (p *GlProject) SetMainProjId(mainProjId string) {
	p.mainProjId = mainProjId
}

func (p *GlProject) addError(err error) {
	if err != nil {
		p.err = err
	}
}
