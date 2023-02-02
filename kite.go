package kite

import "embed"

var (
	Version  = "0.1.0"
	Branch   string
	Revision string

	GoVersion string
	BuildDate string

	PublishAt  = "2021-02-14 13:14"
	UpdatedAt  = "2021-02-14 13:14"
	GithubRepo = "https://github.com/inhere/kite-go"
)

var (
// httproute
)

//go:embed README.md .env.example kite.example.yml config/*.yml
var EmbedFs embed.FS
