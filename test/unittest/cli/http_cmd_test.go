package cli_test

import (
	"testing"

	"github.com/gookit/goutil/testutil/assert"
	"github.com/inhere/kite-go/internal/app"
)

func TestCmd_http_tpl_send(t *testing.T) {
	st := app.Cli.RunLine("http tpl-send -d github --api releases-latest -P -t 2000 -v owner=gookit -v repo=slog")
	assert.Eq(t, st, 0)
}
