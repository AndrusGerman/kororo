package openrouter

import (
	"kororo/internal/adapters/config"
	"kororo/internal/adapters/llm/openai"
	"kororo/internal/core/ports"
)

func New(config *config.Config, model string) ports.LLMAdapter {
	return openai.New(model, config.OPENROUTER_API_KEY(), "https://openrouter.ai/api/v1/v1",
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
