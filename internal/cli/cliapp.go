package cli

import "github.com/gookit/gcli/v3"

// App struct
type App struct {
	*gcli.App
	// FindOsPath find bin in os $PATH
	FindOsPath bool `json:"find_os_path"`
}
