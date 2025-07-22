package kcliapp

import (
	"github.com/gookit/goutil/maputil"
	"github.com/inhere/kite-go/internal/app"
	"github.com/inhere/kite-go/internal/apputil"
)

// kcliapp 用于扩展kite命令，可以通过 yaml, toml 文件定义一个简单的 cli 应用。

// CmdAction struct
type CmdAction struct {
	Name    string   `json:"name"`
	Desc    string   `json:"desc"`
	User    string   `json:"user"`
	Workdir string   `json:"workdir"`
	Cmds    []string `json:"cmds"`
}

// CmdGroup struct
type CmdGroup struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Author  string `json:"author"`
	// Description message
	Description string `json:"description"`
	// DefaultAction on run
	DefaultAction string `json:"default_action"`
	Context       struct {
		Workdir string `json:"workdir"`
	} `json:"context"`
	Actions map[string]*CmdAction `json:"actions"`
}

// ExtraCmds config as command
type ExtraCmds struct {
	maputil.Aliases
	// MetaDir path for load CmdGroup
	MetaDir string `json:"meta_dir"`
	// loaded command groups
	groups map[string]*CmdGroup
}

func (ec *ExtraCmds) Init() error {
	err := app.Cfg().MapOnExists("user_cmds", ec)
	if err != nil {
		return err
	}

	if len(ec.MetaDir) > 0 {
		ec.MetaDir = apputil.ResolvePath(ec.MetaDir)
	}

	return nil
}

func (ec *ExtraCmds) List(name string) {

}

func (ec *ExtraCmds) Info(name string) {

}

func (ec *ExtraCmds) Run(args []string) {

}
