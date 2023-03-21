package devcmd

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/errorx"
)

// ProjectCmd instance
var ProjectCmd = &gcli.Command{
	Name:    "project",
	Aliases: []string{"proj", "projects"},
	Desc:    "simple local projects manage tool commands",
	Config: func(c *gcli.Command) {

	},
	Func: func(c *gcli.Command, _ []string) error {
		return errorx.Raw("TODO")
	},
	Subs: []*gcli.Command{
		ProjectListCmd,
		ProjectAddCmd,
		ProjectRemoveCmd,
		ProjectOpenCmd,
		ProjectRunCmd,
		// TODO list, add, remove, update, run(command), open(by editor)
	},
}

// ProjectListCmd command
var ProjectListCmd = &gcli.Command{
	Name:    "list",
	Aliases: []string{"ls"},
	Desc:    "list all management projects",
	Config: func(c *gcli.Command) {

	},
	Func: func(c *gcli.Command, _ []string) error {
		return errorx.New("TODO")
	},
}

// ProjectAddCmd command
var ProjectAddCmd = &gcli.Command{
	Name:    "add",
	Aliases: []string{"save"},
	Desc:    "add new project to managements",
	Config: func(c *gcli.Command) {

	},
	Func: func(c *gcli.Command, _ []string) error {
		return errorx.New("TODO")
	},
}

// ProjectRemoveCmd command
var ProjectRemoveCmd = &gcli.Command{
	Name:    "del",
	Aliases: []string{"rm", "remove", "delete"},
	Desc:    "remove one or more projects from managements",
	Config: func(c *gcli.Command) {

	},
	Func: func(c *gcli.Command, _ []string) error {
		return errorx.New("TODO")
	},
}

// ProjectOpenCmd command
var ProjectOpenCmd = &gcli.Command{
	Name: "open",
	// Aliases: []string{"rm", "remove", "delete"},
	Desc: "open input project by setting editor",
	Config: func(c *gcli.Command) {

	},
	Func: func(c *gcli.Command, _ []string) error {
		return errorx.New("TODO")
	},
}

// ProjectRunCmd command
var ProjectRunCmd = &gcli.Command{
	Name:    "run",
	Aliases: []string{"exec"},
	Desc:    "execute command on the project dir",
	Config: func(c *gcli.Command) {

	},
	Func: func(c *gcli.Command, _ []string) error {
		return errorx.New("TODO")
	},
}
