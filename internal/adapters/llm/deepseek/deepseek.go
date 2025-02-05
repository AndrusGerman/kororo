package deepseek

import (
	"kororo/internal/adapters/config"
	"kororo/internal/adapters/llm/openai"
	"kororo/internal/core/ports"
)

func New(config *config.Config) ports.LLMAdapter {
	return openai.New("deepseek/deepseek-chat", config.OPENROUTER_API_KEY(), "https://openrouter.ai/api/v1/v1",
		map[string]string{
			"X-Title": config.APP_NAME(),
		},
		map[string]any{
			"provider": map[string]any{
				"sort": "throughput",
			},
		},
	)
}
