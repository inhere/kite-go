module github.com/inherelab/kite

go 1.17

require (
	github.com/fsnotify/fsnotify v1.4.9
	github.com/getkin/kin-openapi v0.22.0
	github.com/go-openapi/spec v0.20.0
	github.com/go-openapi/swag v0.19.12
	github.com/gomarkdown/markdown v0.0.0-20201113031856-722100d81a8e
	github.com/gookit/color v1.5.0
	github.com/gookit/config/v2 v2.0.23
	github.com/gookit/gcli/v3 v3.0.1
	github.com/gookit/gitwrap v0.0.1
	github.com/gookit/goutil v0.4.2
	github.com/gookit/i18n v1.1.3
	github.com/gookit/ini/v2 v2.0.9
	github.com/gookit/rux v1.3.2
	github.com/gookit/slog v0.1.5
	github.com/gookit/view v1.0.2
	github.com/yuin/goldmark v1.3.1
)

require (
	github.com/PuerkitoBio/purell v1.1.1 // indirect
	github.com/PuerkitoBio/urlesc v0.0.0-20170810143723-de5bf2ad4578 // indirect
	github.com/ghodss/yaml v1.0.0 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/jsonreference v0.19.5 // indirect
	github.com/gookit/filter v1.1.2 // indirect
	github.com/gookit/validate v1.2.11 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/mailru/easyjson v0.7.6 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/mapstructure v1.4.1 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/monoculum/formam v0.0.0-20210131081218-41b48e2a724b // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/xo/terminfo v0.0.0-20210125001918-ca9a967f8778 // indirect
	golang.org/x/crypto v0.0.0-20210513164829-c07d793c2f9a // indirect
	golang.org/x/net v0.0.0-20210226172049-e18ecbb05110 // indirect
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c // indirect
	golang.org/x/term v0.0.0-20210503060354-a79de5458b56 // indirect
	golang.org/x/text v0.3.5 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

// for develop
replace github.com/gookit/slog => ../slog

replace github.com/gookit/goutil => ../goutil

replace github.com/gookit/gitwrap => ../gitwrap

replace github.com/gookit/gcli/v3 => ../gcli
