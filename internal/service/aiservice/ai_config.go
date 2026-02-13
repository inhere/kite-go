package aiservice

import (
	"fmt"
	"strings"
)

// ProviderConfig holds the configuration for a single AI provider
type ProviderConfig struct {
	Name        string `json:"name" yaml:"name"`
	BaseURL     string `json:"base_url" yaml:"base_url"`
	APIKey      string `json:"api_key" yaml:"api_key"` // api key, api token
	Description string `json:"description" yaml:"description"`
	Aliases     []string `json:"aliases" yaml:"aliases"`
	// 可用模型列表
	Models   []string `json:"models" yaml:"models"`
	DocsURL  string   `json:"docs_url" yaml:"docs_url"`
	Homepage string   `json:"homepage" yaml:"homepage"`
}

// CCProviderConfig holds the configuration for a single CC provider
type CCProviderConfig struct {
	ProviderConfig `yaml:",squash"`
	// 配置到 cc config 的模型代码
	ModelCode string `json:"model_code" yaml:"model_code"`
	// api key map. key 是自定义名称，value 是 api key
	APIKeys map[string]string `yaml:"api_keys"`
	// 环境变量 需要设置到 cc config
	Envs map[string]string `yaml:"envs"`
}

// GetEnvMaps returns the environment variables for the provider
func (p *CCProviderConfig) GetEnvMaps(keyName, model string) map[string]string {
	if model == "" {
		model = p.ModelCode
	}

	envs := make(map[string]string, len(p.Envs)+2)
	for k, v := range p.Envs {
		envs[strings.ToUpper(k)] = strings.Replace(v, "{model_code}", model, -1)
	}

	envs["ANTHROPIC_BASE_URL"] = p.BaseURL
	apiKey := p.APIKeys[keyName]
	if apiKey == "" {
		apiKey = p.APIKey
	}

	envs["ANTHROPIC_AUTH_TOKEN"] = apiKey
	return envs
}

// Config holds the configuration for the AI service
type Config struct {
	// 默认的模型提供者
	DefaultProvider string `json:"default_provider" yaml:"default_provider"`
	// 默认的模型名称
	DefaultModel    string `json:"default_model" yaml:"default_model"`
	// 支持的模型提供者列表(API场景使用)
	Providers map[string]*ProviderConfig `json:"providers" yaml:"providers"`
	// 提供者别名映射
	ProviderAliases map[string]string `json:"provider_aliases" yaml:"provider_aliases"`
	// 模型别名映射
	ModelAliases map[string]string `json:"model_aliases" yaml:"model_aliases"`
	// 不同场景使用的模型映射
	SceneModels map[string]string `json:"scene_models" yaml:"scene_models"`
	// Claude-code 专用的提供者列表
	CcProviders map[string]*CCProviderConfig `json:"cc_providers" yaml:"cc_providers"`
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

	for name, config := range c.CcProviders {
		config.Name = name
		for _, alias := range config.Aliases {
			c.ProviderAliases[alias] = name
		}
	}
	return nil
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

// ProviderConfig holds the configuration for a single AI provider
func (c *Config) ProviderConfig(name string) (*ProviderConfig, error) {
	name = c.ProviderName(name)
	config, ok := c.Providers[name]
	if !ok {
		return nil, fmt.Errorf("provider %s not found", name)
	}
	return config, nil
}

// IsCCProvider checks if a given name is a valid CC provider
func (c *Config) IsCCProvider(name string) bool {
	name = c.ProviderName(name)
	_, ok := c.CcProviders[name]
	return ok
}

// CCProviderNames returns a list of all available CC provider names
func (c *Config) CCProviderNames() []string {
	names := make([]string, 0, len(c.CcProviders))
	for name := range c.CcProviders {
		names = append(names, name)
	}
	return names
}

// CCProviderConfig holds the configuration for a single CC provider
func (c *Config) CCProviderConfig(name string) (*CCProviderConfig, error) {
	name = c.ProviderName(name)
	config, ok := c.CcProviders[name]
	if !ok {
		return nil, fmt.Errorf("cc-provider %s not found", name)
	}
	config.Name = name
	return config, nil
}
