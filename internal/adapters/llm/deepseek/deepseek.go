package deepseek

import (
	"kororo/internal/adapters/config"
	"kororo/internal/adapters/llm/openai"
	"kororo/internal/core/ports"
)

func New(rest ports.RestAdapter, config *config.Config) ports.LLMAdapter {
	return openai.New(rest, "deepseek-ai/DeepSeek-V3", config.DEEPSEEK_API_KEY(), "https://huggingface.co/api/inference-proxy/together/v1")

}
