package aitool

import "github.com/sashabaranov/go-openai"

type AIConfig struct {
	Provider string `json:"provider" yaml:"provider"` // 默认的模型提供者
	Model    string `json:"model" yaml:"model"`    // 默认的模型名称
	// 支持的模型提供者列表
	Providers []string `json:"providers" yaml:"providers"`
	// 提供者别名映射
	ProviderAliases map[string]string `json:"provider_aliases" yaml:"provider_aliases"`
	// 模型别名映射
	ModelAliases map[string]string `json:"model_aliases" yaml:"model_aliases"`
	// 不同场景使用的模型映射
	SceneModels map[string]string `json:"scene_models" yaml:"scene_models"`
}

type AIClient struct {
	cfg *AIConfig
	*openai.Client
}

var std = &AIClient{}

func Client() *AIClient {
	return std
}

func Init(cfg *AIConfig) {
	std.cfg = cfg
}
