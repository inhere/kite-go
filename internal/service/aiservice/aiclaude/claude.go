package aiclaude

// ClaudeRuntimeConfig represents the Claude configuration file format
//
//  - Linux, Mac: ~/.claude/config.json
//  - Windows: %USERPROFILE%\.claude\settings.json
type ClaudeRuntimeConfig struct {
	Env map[string]string `json:"env,omitempty"`
	// IncludeCoAuthoredBy indicates whether to include co-authored-by in the commit message
	IncludeCoAuthoredBy bool `json:"includeCoAuthoredBy"`
	// enabledPlugins map
	EnabledPlugins map[string]string `json:"enabledPlugins,omitempty"`
	// statusLine string map
	StatusLine map[string]string `json:"statusLine,omitempty"`
}

// Save saves the configuration to the file
func (c *ClaudeRuntimeConfig) Save() error {
	return WriteUserConfig(c)
}
