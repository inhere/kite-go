package taskx

import "github.com/gookit/gcli/v3"

// TaskManageCmd instance
var TaskManageCmd = &gcli.Command{
	Name: "task",
	Desc: "Task manage tools command",
	Subs: []*gcli.Command{
		TaskList,
		TaskInfo,
		TaskRun,
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
