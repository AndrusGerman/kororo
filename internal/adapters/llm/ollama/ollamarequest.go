package ollama

type ollamaRequest struct {
	Model    string     `json:"model"`
	Stream   bool       `json:"stream"`
	Messages []*message `json:"messages"`
}

func newOllamaRequest(model string, messages []*message, stream bool) *ollamaRequest {
	return &ollamaRequest{
		Model:    model,
		Messages: messages,
		Stream:   stream,
	}
}
