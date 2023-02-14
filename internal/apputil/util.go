package apputil

import (
	"strings"

	"github.com/gookit/goutil/maputil"
	"github.com/inhere/kite/internal/app"
)

// CmdConfigKey build
func CmdConfigKey(nodes ...string) string {
	return "cmd_" + strings.Join(nodes, "_")
}

// CmdConfigData find
func CmdConfigData(nodes ...string) maputil.Data {
	return app.Cfg().SubDataMap(CmdConfigKey(nodes...))
}
