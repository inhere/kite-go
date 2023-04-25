package appcmd

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/gflag"
	"github.com/gookit/goutil/byteutil"
	"github.com/gookit/slog"
	"github.com/inhere/kite-go"
	"github.com/inhere/kite-go/internal/app"
	"github.com/inhere/kite-go/internal/apputil"
)

// ReadmeCmd instance
var ReadmeCmd = &gcli.Command{
	Name:    "readme",
	Aliases: []string{"doc", "docs"},
	Desc:    "show readme docs for kite app",
	Config: func(c *gcli.Command) {

	},
	Func: func(c *gcli.Command, _ []string) error {
		str := byteutil.SafeString(kite.EmbedFs.ReadFile("README.md"))
		// fmt.Println(byteutil.SafeString(kite.EmbedFs.ReadFile("README.md")))

		return apputil.RenderContents(str, "markdown", "github")
	},
}

var lwOpts = struct {
	Level string `flag:"set log level;false;info;l"`
}{}

// LogWriteCmd instance
var LogWriteCmd = &gcli.Command{
	Name:    "logw",
	Aliases: []string{"log"},
	Desc:    "write a log message to the kite app logs file",
	Config: func(c *gcli.Command) {
		c.MustFromStruct(&lwOpts, gflag.TagRuleSimple)
		c.AddArg("message", "the log message", true)
	},
	Func: func(c *gcli.Command, _ []string) error {
		app.L.Log(slog.LevelByName(lwOpts.Level), c.Arg("message").String())
		return nil
	},
}
