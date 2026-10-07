package gitcmd

import (
	"strings"

	"github.com/gookit/goutil/strutil"
)

// VersionBumpLevel the level to bump a version.
type VersionBumpLevel int

const (
	// BumpMajor bump the major version. eg: v1.2.3 -> v2.0.0
	BumpMajor VersionBumpLevel = iota
	// BumpMinor bump the minor version. eg: v1.2.3 -> v1.3.0
	BumpMinor
	// BumpPatch bump the patch version. eg: v1.2.3 -> v1.2.4
	BumpPatch
)

// ParseBumpLevel parse the bump level from a string.
// supported: major|maj|1, minor|min|2, patch|fix|3. empty will use patch.
func ParseBumpLevel(level string) (VersionBumpLevel, bool) {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "", "patch", "fix", "p", "3":
		return BumpPatch, true
	case "minor", "min", "2":
		return BumpMinor, true
	case "major", "maj", "1":
		return BumpMajor, true
	}
	return BumpPatch, false
}

// BumpVersion bump the version to next by the level. the lower levels will be reset to 0.
// eg: BumpVersion("v1.2.3", BumpMinor) -> "1.3.0"
func BumpVersion(ver string, level VersionBumpLevel) string {
	ver = strings.TrimLeft(strings.TrimSpace(ver), "vV")
	// drop the pre-release suffix. eg: 1.2.3-beta -> 1.2.3
	if idx := strings.IndexByte(ver, '-'); idx >= 0 {
		ver = ver[:idx]
	}
	if ver == "" {
		ver = "0.0.0"
	}

	nodes := strings.Split(ver, ".")
	for len(nodes) < 3 {
		nodes = append(nodes, "0")
	}

	idx := int(level)
	if idx < 0 || idx >= len(nodes) {
		idx = len(nodes) - 1
	}

	num, err := strutil.ToInt(nodes[idx])
	if err != nil {
		num = 0
	}
	nodes[idx] = strutil.SafeString(num + 1)

	for i := idx + 1; i < len(nodes); i++ {
		nodes[i] = "0"
	}
	return strings.Join(nodes, ".")
}
