package gituse

import (
	"errors"

	"github.com/gookit/gcli/v3"
)

// OpenRemoteRepo address
var OpenRemoteRepo = &gcli.Command{
	Name: "open",
	Desc: "open the git remote repo address",
	Func: func(c *gcli.Command, args []string) error {
		return errors.New("TODO")
	},
}
