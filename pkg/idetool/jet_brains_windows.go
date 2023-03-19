package idetool

import "github.com/gookit/goutil/sysutil"

// ProfileDir for ide
func (j *JetBrains) ProfileDir() string {
	if j.profileDir == "" {
		j.profileDir = sysutil.ExpandHome("~/AppData/Roaming/JetBrains")
	}
	return j.profileDir
}
