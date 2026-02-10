package aiservice

import "fmt"

// ProviderConfig holds the configuration for a single AI provider
type ProviderConfig struct {
	Name        string `json:"name" yaml:"name"`
	BaseURL     string `json:"base_url" yaml:"base_url"`
	APIKey      string `json:"api_key" yaml:"api_key"` // api key, api token
	Description string `json:"description" yaml:"description"`
	Aliases     []string `json:"aliases" yaml:"aliases"`
}

// GenShellEnv generates the environment variables
func (p *ProviderConfig) GenShellEnv(shell string) string {
	shellTpl  := `# Claude API Configuration (use %s)
export ANTHROPIC_BASE_URL="%s"
export ANTHROPIC_AUTH_TOKEN="%s"
`

	if shell == "pwsh" || shell == "powershell" {
		shellTpl = `# Claude API Configuration (use pwsh)
$env:ANTHROPIC_BASE_URL="%s"
$env:ANTHROPIC_AUTH_TOKEN="%s"
`
	}
	return fmt.Sprintf(shellTpl, p.Name, p.BaseURL, p.APIKey)
}

// Config holds the configuration for the AI service
type Config struct {
	// 默认的模型提供者
	DefaultProvider string `json:"default_provider" yaml:"default_provider"`
	// 默认的模型名称
	DefaultModel    string `json:"default_model" yaml:"default_model"`
	// 支持的模型提供者列表
	Providers map[string]ProviderConfig `json:"providers" yaml:"providers"`
	// 提供者别名映射
	ProviderAliases map[string]string `json:"provider_aliases" yaml:"provider_aliases"`
	// 模型别名映射
	ModelAliases map[string]string `json:"model_aliases" yaml:"model_aliases"`
	// 不同场景使用的模型映射
	SceneModels map[string]string `json:"scene_models" yaml:"scene_models"`
}

// Init 初始化 config 部分信息
func (c *Config) Init() error {
	// 模型别名映射
	if c.ModelAliases == nil {
		c.ModelAliases = make(map[string]string)
	}
	// 模型别名映射
	if c.SceneModels == nil {
		c.SceneModels = make(map[string]string)
	}

	for name, config := range c.Providers {
		config.Name = name
		for _, alias := range config.Aliases {
			c.ProviderAliases[alias] = name
		}
	}

	return nil
}

// ProviderConfig holds the configuration for a single AI provider
func (c *Config) ProviderConfig(name string) (*ProviderConfig, error) {
	name = c.ProviderName(name)
	config, ok := c.Providers[name]
	if !ok {
		return nil, fmt.Errorf("provider %s not found", name)
	}
	return &config, nil
}

// ModelName returns the real name of a model
func (c *Config) ModelName(name string) string {
	if realName, ok := c.ModelAliases[name]; ok {
		return realName
	}
	return name
}

// ProviderName returns the real name of a provider
func (c *Config) ProviderName(name string) string {
	if realName, ok := c.ProviderAliases[name]; ok {
		return realName
	}
	return name
}

// IsProvider checks if a given name is a valid provider
func (c *Config) IsProvider(name string) bool {
	name = c.ProviderName(name)
	_, ok := c.Providers[name]
	return ok
}

// ProviderNames returns a list of all available provider names
func (c *Config) ProviderNames() []string {
	names := make([]string, 0, len(c.Providers))
	for name := range c.Providers {
		names = append(names, name)
	}
	return names
}
