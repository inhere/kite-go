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

//go:embed README.md .example.env kite.example.yml config
var EmbedFs embed.FS

// Banner text
// from http://patorjk.com/software/taag/#p=testall&f=Graffiti&t=Kite
// font: Doom,Graffiti,Isometric1 - Isometric3, Ogre, Slant
var Banner = `
GoVersion: {{goVersion}}
BuildDate: {{buildDate}}
 __  __     __     ______   ______
/\ \/ /    /\ \   /\__  _\ /\  ___\
\ \  _"-.  \ \ \  \/_/\ \/ \ \  __\
 \ \_\ \_\  \ \_\    \ \_\  \ \_____\
  \/_/\/_/   \/_/     \/_/   \/_____/

Auther  : https://github.com/inhere
Homepage: https://github.com/inhere/kite-go
`
