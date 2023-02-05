package kiteext

// PluginRunner struct
type PluginRunner struct {
	// DenyNames deny plugin bin names
	DenyNames []string `json:"deny_names"`
	// PluginDirs plugin bin search dirs
	PluginDirs []string `json:"plugin_dirs"`
	ConfigFile []string `json:"config_file"`
}
