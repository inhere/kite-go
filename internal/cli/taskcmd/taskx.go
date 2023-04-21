package taskcmd

import "github.com/gookit/gcli/v3"

// TaskManageCmd instance
var TaskManageCmd = &gcli.Command{
	Name:    "task",
	Desc:    "Task/Script run and management command",
	Aliases: []string{"scripts", "script"},
	Subs: []*gcli.Command{
		TaskList,
		TaskInfo,
		TaskRun,
	},
	Config: func(c *gcli.Command) {

	},
	Func: func(c *gcli.Command, args []string) error {
		return nil
	},
}

var TaskList = &gcli.Command{
	Name:    "list",
	Desc:    "list all Tasks",
	Aliases: []string{"ls", "l"},
}

var TaskInfo = &gcli.Command{
	Name: "info",
	Desc: "show an Task information",
}

var TaskRun = &gcli.Command{
	Name: "run",
	Desc: "run an Task",
}
