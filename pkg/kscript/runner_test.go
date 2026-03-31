package kscript

import "testing"

func TestIfExpr(t *testing.T) {

	st := &ScriptTask{
		If: "",
	}

	st.resolveIfExpr(map[string]any{
		"make": true,
	})
}
