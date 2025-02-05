package ollama

import (
	"kororo/internal/adapters/llm/openai"
	"kororo/internal/core/ports"
)

func New() ports.LLMAdapter {
	return openai.New("gemma2:27b-instruct-q4_K_S", "", "http://localhost:11434/v1/v1", map[string]string{}, map[string]interface{}{})
}
