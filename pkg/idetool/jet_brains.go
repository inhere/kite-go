package idetool

import (
	"io/fs"

	"github.com/gookit/goutil/fsutil"
)

const (
	IDEIntellij = "Intellij"
	IDEGoLand   = "GoLand"
	IDEPhpStorm = "PhpStorm"
	IDEWebStorm = "WebStorm"
)

// JetBrains struct
type JetBrains struct {
	installDir string
	profileDir string
}

// NewJetBrains instance
func NewJetBrains() *JetBrains {
	return &JetBrains{}
}

// HasToolbox on local
func (j *JetBrains) HasToolbox() bool {
	return fsutil.IsDir(j.ToolboxDir())
}

// ToolboxDir on local
func (j *JetBrains) ToolboxDir() string {
	return j.ProfileDir() + "/Toolbox"
}

// Installed tools on local
//
// macOS: ~/Library/Application Support/JetBrains/GoLand2022.3
func (j *JetBrains) Installed() map[string]string {
	mp := make(map[string]string)
	// ~/Library/Application\ Support/JetBrains/Toolbox/apps/
	if j.HasToolbox() {
		_ = fsutil.FindInDir(j.ProfileDir(), func(fPath string, ent fs.DirEntry) error {
			mp[ent.Name()] = fPath
			return nil
		}, fsutil.OnlyFindDir)
		return mp
	}

	_ = fsutil.FindInDir(j.ProfileDir(), func(fPath string, ent fs.DirEntry) error {
		mp[ent.Name()] = fPath
		return nil
	}, fsutil.OnlyFindDir, fsutil.ExcludeSuffix("-backup", "Options", "Policy"))
	return mp
}
