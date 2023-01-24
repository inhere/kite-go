package gitlab

import (
	"net/url"
	"strings"

	"github.com/gookit/goutil/strutil"
)

// PRLinkBuilder struct
type PRLinkBuilder struct {
	HostUrl     string
	RepoPath    string
	QueryString string
}

// PRLinkQuery struct
type PRLinkQuery struct {
	RepoPath string
	// SourceProjectId int or urlencode(group/name)
	SourceProjectId string `json:"source_project_id"`
	SourceBranch    string `json:"source_branch"`
	// TargetProjectId int or urlencode(group/name)
	TargetProjectId string `json:"target_project_id"`
	TargetBranch    string `json:"target_branch"`
}

// NewPRLinkQuery instance
func NewPRLinkQuery(srcPid, srcBr, dstPid, dstBr string) *PRLinkQuery {
	return &PRLinkQuery{
		SourceProjectId: srcPid,
		SourceBranch:    srcBr,
		TargetProjectId: strutil.OrElse(srcPid, dstPid),
		TargetBranch:    strutil.OrElse(srcBr, dstBr),
	}
}

// QueryString build query string.
func (b *PRLinkQuery) QueryString() string {
	qry := url.Values{}
	qry.Add("utf8", "✓")

	qry.Add("merge_request[source_project_id]", b.SourceProjectId)
	qry.Add("merge_request[source_branch]", b.SourceBranch)
	qry.Add("merge_request[target_project_id]", b.TargetProjectId)
	qry.Add("merge_request[target_branch]", b.TargetBranch)

	return qry.Encode()
}

// BuildURL pr link url
func (b *PRLinkQuery) BuildURL(hostUrl string) string {
	var sb strings.Builder
	sb.Grow(320)

	sb.WriteString(hostUrl)
	sb.WriteByte('/')
	sb.WriteString(b.RepoPath)
	sb.WriteByte('?')
	sb.WriteString(b.QueryString())

	return sb.String()
}
