package cmdbiz

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/arrutil"
	"github.com/gookit/goutil/strutil"
	"github.com/inhere/kite-go/internal/apputil"
)

// ProxyCmdConf struct
type ProxyCmdConf struct {
	// CommandIds eg: [git, git:clone, git:tag:*]
	CommandIds []string `json:"command_ids"`
	// GroupLimits eg: {github: {acp, update}}
	GroupLimits map[string][]string `json:"group_limits"`
}

// ProxyCC instance
var ProxyCC = &ProxyCmdConf{}

// AutoSetByName handle
func (pcc *ProxyCmdConf) AutoSetByName(group, sub string, args []string) {
	if pcc.IsMatchName(group, sub, args) {
		apputil.ApplyProxyEnv()
	}
}

// AutoSetByCmd handle
func (pcc *ProxyCmdConf) AutoSetByCmd(c *gcli.Command) {
	if pcc.IsMatchCmd(c) {
		apputil.ApplyProxyEnv()
	}
}

// IsMatchCmd by config
func (pcc *ProxyCmdConf) IsMatchCmd(c *gcli.Command) bool {
	cmdId := c.ID()
	if arrutil.StringsHas(pcc.CommandIds, cmdId) {
		return true
	}

	if names := c.PathNames(); len(names) > 1 {
		group, sub := names[0], strutil.JoinList(gcli.CommandSep, names[1:])
		if subs, ok := pcc.GroupLimits[group]; ok {
			return arrutil.StringsHas(subs, sub)
		}
	}

	return false
}

// IsMatchName by config
func (pcc *ProxyCmdConf) IsMatchName(group, sub string, args []string) bool {
	// collect all arguments
	nodes := []string{group, sub}
	for _, arg := range args {
		if len(arg) == 0 || arg[0] == '-' {
			break
		}
		nodes = append(nodes, arg)
	}

	// match command ids
	cmdId := strutil.Join(gcli.CommandSep, group, sub)
	cmdId2 := strutil.JoinList(gcli.CommandSep, nodes)
	for _, idStr := range pcc.CommandIds {
		if cmdId == idStr || cmdId2 == idStr {
			return true
		}

		if strutil.MatchNodePath(idStr, cmdId2, gcli.CommandSep) {
			return true
		}
	}

	if subs, ok := pcc.GroupLimits[group]; ok {
		return arrutil.StringsHas(subs, sub)
	}
	return false
}
