package extcmd

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/maputil"
	"github.com/inhere/kite-go/internal/app"
	"github.com/inhere/kite-go/internal/apputil"
)

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

var ueOpts = struct {
	list bool
	info string
}{}

// UserExtCmd instance
var UserExtCmd = &gcli.Command{
	Name:    "xcmd",
	Aliases: []string{"ucmd"},
	Desc:    "manage and execute user extension commands",
	Subs: []*gcli.Command{
		ListUserCmd,
	},
	Config: func(c *gcli.Command) {
		c.BoolOpt2(&ueOpts.list, "list, ls", "display all user extra commands")
		c.StrOpt2(&ueOpts.info, "info, show, i", "display info for the input command")

		c.AddArg("command", "input command for execute")
	},
	Func: func(c *gcli.Command, args []string) error {

		name := c.Arg("command").String()
		if len(name) == 0 {
			return c.NewErr("please input command name for handle")
		}

		return nil
	},
}

var ListUserCmd = &gcli.Command{
	Name:    "list",
	Desc:    "list all added user commands",
	Aliases: []string{"ls", "l"},
}
