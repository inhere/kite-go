module github.com/inherelab/kite

go 1.16

require (
	github.com/fsnotify/fsnotify v1.5.4
	github.com/getkin/kin-openapi v0.22.0
	github.com/go-openapi/spec v0.20.0
	github.com/go-openapi/swag v0.19.12
	github.com/gomarkdown/markdown v0.0.0-20201113031856-722100d81a8e
	github.com/gookit/color v1.5.0
	github.com/gookit/config/v2 v2.0.27
	github.com/gookit/gcli/v3 v3.0.1
	github.com/gookit/gitw v0.2.1
	github.com/gookit/goutil v0.5.1
	github.com/gookit/i18n v1.1.3
	github.com/gookit/ini/v2 v2.0.11
	github.com/gookit/rux v1.3.2
	github.com/gookit/slog v0.2.1
	github.com/gookit/view v1.0.2
	github.com/yuin/goldmark v1.3.1
)

// for develop
replace github.com/gookit/slog => ../slog

replace github.com/gookit/goutil => ../goutil

replace github.com/gookit/gitw => ../gitw

replace github.com/gookit/gcli/v3 => ../gcli
