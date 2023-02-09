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

//go:embed README.md .example.env kite.example.yml config/*.yml
var EmbedFs embed.FS

// Banner text
// from http://patorjk.com/software/taag/#p=testall&f=Graffiti&t=Kite
// font: Doom,Graffiti,Isometric1 - Isometric3, Ogre, Slant
var Banner = `
      ___                                   ___
     /__/|        ___           ___        /  /\
    |  |:|       /  /\         /  /\      /  /:/_
    |  |:|      /  /:/        /  /:/     /  /:/ /\
  __|  |:|     /__/::\       /  /:/     /  /:/ /:/_
 /__/\_|:|____ \__\/\:\__   /  /::\    /__/:/ /:/ /\
 \  \:\/:::::/    \  \:\/\ /__/:/\:\   \  \:\/:/ /:/
  \  \::/~~~~      \__\::/ \__\/  \:\   \  \::/ /:/
   \  \:\          /__/:/       \  \:\   \  \:\/:/
    \  \:\         \__\/         \__\/    \  \::/
     \__\/                                 \__\/
`
