package gencmd

import "github.com/gookit/gcli/v3"

type projectOpt struct {
	Dir  string `flag:"desc=the directory for create project"`
	Tpl  string `flag:"desc=the template project path, allow: local path, git url"`
	Conf string `flag:"desc=the config file path;shorts=c,config"`
}

// NewProjectCmd create a new project command
func NewProjectCmd() *gcli.Command {
	opt := projectOpt{}

	return &gcli.Command{
		Name:    "project",
		Desc:    "create a new project",
		Aliases: []string{"proj", "prj"},
		Config: func(c *gcli.Command) {
			c.MustFromStruct(&opt)
		},
		Func: func(c *gcli.Command, args []string) error {
			return c.NewErr("TODO")
		},
	}
}
