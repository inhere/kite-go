package bootstrap

import "github.com/inherelab/kite/app"

func test() {
	app.Cfg().MapOnExists()
}
