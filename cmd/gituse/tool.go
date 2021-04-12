package gituse

import "github.com/gookit/gcli/v3"

var BatchPull = &gcli.Command{
	Name: "bpull",
	Desc: "batch pull multi git directory by `git pull`",
	Aliases: []string{"bp", "batch-pull"},
}
