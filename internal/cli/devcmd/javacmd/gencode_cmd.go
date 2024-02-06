package javacmd

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/errorx"
)

func NewGenCodeCmd() *gcli.Command {
	return &gcli.Command{
		Name: "gen-code",
		Desc: "generate code for java project",
		Func: func(c *gcli.Command, _ []string) error {
			return errorx.New("TODO")
		},
	}
}

func NewGenDtoCmd() *gcli.Command {
	return &gcli.Command{
		Name: "gen-code",
		Desc: "generate code for java project",
		Func: func(c *gcli.Command, _ []string) error {
			return errorx.New("TODO")
		},
	}
}
