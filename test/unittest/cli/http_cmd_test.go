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

func TestCmd_http_send(t *testing.T) {
	st := app.Cli.RunLine(`http send -j -m POST -H 'Content-Type: application/json' -d '{"project":"test-demo","tag": "1.1.0"}' http://172.20.0.6:9091/tag/add`)
	assert.Eq(t, st, 0)
}
