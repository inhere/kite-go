package idetool

import "github.com/gookit/goutil/sysutil"

// ProfileDir for ide
func (j *JetBrains) ProfileDir() string {
	// ~/Library/Application Support/JetBrains/GoLand2022.3
	if j.profileDir == "" {
		j.profileDir = sysutil.ExpandHome("~/Library/Application Support/JetBrains")
	}
	return j.profileDir
}
