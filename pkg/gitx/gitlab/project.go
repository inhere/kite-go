package gitlab

import (
	"net/url"
	"strings"

	"github.com/gookit/gitw"
	"github.com/gookit/goutil"
	"github.com/gookit/goutil/errorx"
	"github.com/inhere/kite-go/pkg/gitx"
)

// GlProject struct
type GlProject struct {
	*GitLab

	lp *gitx.GitLoc

	// remote info
	mainRemoteInfo *gitw.RemoteInfo
	forkRemoteInfo *gitw.RemoteInfo

	// other
	mainProjId, forkProjId string
}

// NewGlProject instance
func NewGlProject(dir string, gl *GitLab) *GlProject {
	return &GlProject{
		GitLab: gl,
		lp:     gl.LoadRepo(dir),
	}
}

// GitLoc instance
func (p *GlProject) GitLoc() *gitx.GitLoc {
	return p.lp
}

func (p *GlProject) CheckRemote() error {
	if err := p.CheckDefaultRemote(); err != nil {
		return err
	}
	return p.CheckSourceRemote()
}

func (p *GlProject) CheckDefaultRemote() error {
	if !p.lp.HasRemote(p.DefaultRemote) {
		return errorx.Newf("the fork remote '%s' is not found(config:gitlab.default_remote)", p.DefaultRemote)
	}
	return nil
}

func (p *GlProject) CheckSourceRemote() error {
	if !p.lp.HasRemote(p.SourceRemote) {
		return errorx.Newf("the main remote '%s' is not found(config:gitlab.source_remote)", p.SourceRemote)
	}
	return nil
}

func (p *GlProject) ForkRmtInfo() *gitw.RemoteInfo {
	if p.forkRemoteInfo == nil {
		goutil.PanicErr(p.CheckDefaultRemote())
		p.forkRemoteInfo = p.lp.RemoteInfo(p.DefaultRemote)
	}
	return p.forkRemoteInfo
}

func (p *GlProject) MainRmtInfo() *gitw.RemoteInfo {
	if p.mainRemoteInfo == nil {
		goutil.PanicErr(p.CheckDefaultRemote())
		p.mainRemoteInfo = p.lp.RemoteInfo(p.SourceRemote)
	}
	return p.mainRemoteInfo
}

// ResolveBranch name
func (p *GlProject) ResolveBranch(brName string) (string, bool) {
	switch strings.ToUpper(brName) {
	case "", "@", "H", "HEAD":
		return p.lp.CurBranchName(), true
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
