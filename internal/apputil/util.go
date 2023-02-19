package apputil

import (
	"strings"

	"github.com/gookit/goutil/maputil"
	"github.com/inhere/kite/internal/app"
	"github.com/inhere/kite/pkg/kiteext"
)

// CmdConfigKey build
func CmdConfigKey(nodes ...string) string {
	return "cmd_" + strings.Join(nodes, "_")
}

// CmdConfigData find
func CmdConfigData(nodes ...string) maputil.Data {
	return app.Cfg().SubDataMap(CmdConfigKey(nodes...))
}

// ReadSource string data.
func ReadSource(s string) string {
	return kiteext.NewSourceReader(s).ReadString()
}
