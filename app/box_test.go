package app_test

import (
	"testing"

	"github.com/gookit/goutil/dump"
	"github.com/gookit/goutil/testutil/assert"
	"github.com/inherelab/kite/app"
)

type testObj struct {
	Name string
}

func TestGet(t *testing.T) {
	app.Add("obj1", &testObj{Name: "inhere"})

	obj1 := app.Get[*testObj]("obj1")
	dump.P(obj1)

	obj2, ok := app.Lookup[*testObj]("obj2")
	assert.False(t, ok)
	assert.Nil(t, obj2)
}
