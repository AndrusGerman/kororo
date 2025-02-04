package deepseek

import (
	"kororo/internal/adapters/config"
	"kororo/internal/adapters/llm/openai"
	"kororo/internal/core/ports"
)

func New(config *config.Config) ports.LLMAdapter {
	return openai.New("deepseek-ai/DeepSeek-V3", config.HUGGINGFACE_API_KEY(), "https://huggingface.co/api/inference-proxy/together/v1")

}
