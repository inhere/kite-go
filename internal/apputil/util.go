package apputil

import (
	"strings"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/gitw"
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

// CmdConfigKey build. eg: ("git", "update") => "cmd_git_update"
func CmdConfigKey(nodes ...string) string {
	return "cmd_" + strings.Join(nodes, "_")
}

// ReadSource string data.
func ReadSource(s string) string {
	return kiteext.NewSourceReader(s).ReadString()
}

// GitCfgByCmdID get
func GitCfgByCmdID(c *gcli.Command) *gitx.Config {
	id := c.ID()
	if strings.Contains(id, gitw.TypeGitHub) {
		return app.Ghub().Config
	}

	if strings.Contains(id, gitw.TypeGitlab) {
		return app.Glab().Config
	}
	return app.Gitx()
}
