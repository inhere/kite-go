package kiteext

import (
	"github.com/gookit/goutil/errorx"
	"github.com/gookit/goutil/maputil"
)

// KiteAliasRun struct
type KiteAliasRun struct {
	Aliases maputil.Aliases `json:"aliases"`
}

func (kar *KiteAliasRun) IsAlias(name string) bool {
	return kar.Aliases.HasAlias(name)
}

func (kar *KiteAliasRun) Run(name string, args []string) error {
	if !kar.Aliases.HasAlias(name) {
		return errorx.Newf("kite alias command %q is not found", name)
	}

	return nil
}
