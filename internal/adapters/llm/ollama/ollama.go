package ollama

import (
	"kororo/internal/adapters/llm/openai"
	"kororo/internal/core/ports"
)

func New(rest ports.RestAdapter) ports.LLMAdapter {
	return openai.New(rest, "mistral:latest", "", "http://localhost:11434/api/chat")
}
