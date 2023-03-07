package apputil

import (
	"strings"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/gitw"
	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/maputil"
	"github.com/inhere/kite/internal/app"
	"github.com/inhere/kite/pkg/gitx"
	"github.com/inhere/kite/pkg/kiteext"
)

// CmdConfigData find.
//
// eg:
//
//	CmdConfigData("git", "update") => read config by key: cmd_git_update
func CmdConfigData(nodes ...string) maputil.Data {
	return app.Cfg().SubDataMap(CmdConfigKey(nodes...))
}

// CmdConfigData2 find.
//
// eg:
//
//	CmdConfigData2("git", "update") => read config by key: cmd_git_update
func CmdConfigData2(c *gcli.Command) maputil.Data {
	return app.Cfg().SubDataMap(CmdConfigKey(c.PathNames()...))
}

// CmdConfigKey build. eg: ("git", "update") => "cmd_git_update"
func CmdConfigKey(nodes ...string) string {
	return "cmd_" + strings.Join(nodes, "_")
}

// ReadSource string data.
func ReadSource(s string) (string, error) {
	return kiteext.
		NewSourceReader(s, kiteext.FallbackStdin(), kiteext.WithTrimSpace()).
		TryReadString()
}

// GitCfgByCmdID get
func GitCfgByCmdID(c *gcli.Command) (cfg *gitx.Config) {
	id := c.ID()

	if strings.Contains(id, gitw.TypeGitHub) {
		cfg = app.Ghub().Config
	} else if strings.Contains(id, gitw.TypeGitlab) {
		cfg = app.Glab().Config
	} else {
		cfg = app.Gitx()
	}

	c.Infof("TIP: auto select git config type: %s(by cmd ID: %s)\n", cfg.HostType, c.ID())
	return cfg
}

// ResolvePath for input path
func ResolvePath(path string) string {
	path = app.Vars.Replace(path)

	if app.IsAliasPath(path) {
		return app.App().ResolvePath(path)
	}

	if fsutil.IsAbsPath(path) {
		return path
	}
	return app.App().PathBuild(path)
}

// ResolveSep char
func ResolveSep(sep string) string {
	switch sep {
	case "SPACE":
		return " "
	case "NL", "NEWLINE":
		return "\n"
	case "TAB":
		return "\t"
	default:
		return sep
	}
}
