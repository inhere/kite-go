package gitx

import (
	"regexp"

	"github.com/gookit/goutil/strutil"
)

// BranchMatcher interface
type BranchMatcher interface {
	// Match branch name, no remote prefix
	Match(branch string) bool
}

// GlobMatch handle glob matching
type GlobMatch struct {
	pattern string
}

// Match branch name by glob pattern
func (g *GlobMatch) Match(branch string) bool {
	return strutil.GlobMatch(g.pattern, branch)
}

// RegexMatch handle regex matching
type RegexMatch struct {
	pattern string
	regex   *regexp.Regexp
}

// Match branch name by regex pattern
func (r *RegexMatch) Match(branch string) bool {
	return r.regex.MatchString(branch)
}

// NewBranchMatcher create a new branch matcher
func NewBranchMatcher(pattern string, regex bool) BranchMatcher {
	if regex {
		return &RegexMatch{
			pattern: pattern,
			regex:   regexp.MustCompile(pattern),
		}
	}
	return &GlobMatch{pattern: pattern}
}
