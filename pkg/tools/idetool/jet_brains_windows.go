package idetool

import (
	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/sysutil"
)

// ProfileDir for ide
func (j *JetBrains) ProfileDir(subPath ...string) string {
	if j.profileDir == "" {
		j.profileDir = sysutil.ExpandHome("~/AppData/Roaming/JetBrains")
	}

	if len(subPath) > 0 {
		return fsutil.JoinSubPaths(j.profileDir, subPath...)
	}
	return j.profileDir
}
