package main

import (
	"github.com/gookit/gcli/v2"
	"github.com/gookit/gcli/v2/builtin"
	"github.com/gookit/kite/cmd/mkdown"

	// "github.com/gookit/gcli/v2/builtin/reverseproxy"
	"runtime"
)

// local run:
// 	go run ./_examples/cliapp.go
// 	go build ./_examples/cliapp.go && ./cliapp
//
// run on windows(cmd, powerShell):
// 	go build ./_examples/cliapp.go; ./cliapp
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	app := gcli.NewApp(func(app *gcli.App) {
		app.Version = "1.0.6"
		app.Description = "this is my cli application"
		app.On(gcli.EvtAppInit, func(data ...interface{}) {
			// do something...
			// fmt.Println("init app")
		})

		// app.SetVerbose(gcli.VerbDebug)
		// app.DefaultCommand("example")
		app.Logo.Text = `
   __________          __
  /  _/_  __/__  ___  / /
 _/ /  / / / _ \/ _ \/ /
/___/ /_/  \___/\___/_/
`
	})

	// app.Strict = true
	// app.Add(filewatcher.FileWatcher(nil))
	// app.Add(reverseproxy.ReverseProxyCommand())

	app.Add(&gcli.Command{
		Name:    "test",
		Aliases: []string{"ts"},
		UseFor:  "this is a description <info>message</> for command {$cmd}",
		Func: func(cmd *gcli.Command, args []string) error {
			gcli.Print("hello, in the test command\n")
			return nil
		},
	})

	app.Add(builtin.GenAutoComplete())
	app.Add(mkdown.ConvertMD2html())

	// running
	app.Run()
}
