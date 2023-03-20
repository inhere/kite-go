package appcmd

import (
	"fmt"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/fsutil"
	"github.com/inhere/kite-go/internal/app"
	"github.com/inhere/kite-go/internal/appconst"
)

// AppCheckCmd command
var AppCheckCmd = &gcli.Command{
	Name: "check",
	Desc: "check the kite app runtime env information",
	Func: func(c *gcli.Command, args []string) error {

		baseDir := app.App().BaseDir
		confDir := app.App().ConfigDir
		dotenvFile := fsutil.JoinPaths(baseDir, appconst.DotEnvFileName)

		infos := []struct {
			title  string
			answer string
		}{
			{
				"Kite base data dir: " + baseDir,
				fmt.Sprintf("Exists: %v", fsutil.IsDir(baseDir)),
			},
			{
				"Kite dotenv file: " + dotenvFile,
				fmt.Sprintf("Exists: %v", fsutil.IsFile(dotenvFile)),
			},
			{
				"Kite config directory: " + confDir,
				fmt.Sprintf("Exists: %v", fsutil.IsDir(confDir)),
			},
		}

		// fsutil.FileTree()
		for _, info := range infos {
			c.Println("-", info.title)
			c.Infoln(">", info.answer)
		}

		return nil
	},
}
