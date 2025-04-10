package apputil

import (
	"github.com/gookit/goutil/envutil"
	"github.com/sashabaranov/go-openai"
)

// NewOpenaiClient create a new OpenAI client
func NewOpenaiClient() *openai.Client {
	openaiConfig := openai.DefaultConfig(envutil.MustGet("OPENAI_API_KEY"))
	envutil.OnExist("OPENAI_BASE_URL", func(val string) {
		openaiConfig.BaseURL = val
	})
	envutil.OnExist("LLM_SERVICE_TYPE", func(val string) {
		openaiConfig.APIType = openai.APIType(val)
	})

	return openai.NewClientWithConfig(openaiConfig)
}
