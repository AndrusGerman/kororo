package huggingface

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func newMessage(role string, content string) *message {
	return &message{
		Role:    role,
		Content: content,
	}
}

type huggingfaceRequest struct {
	Model     string     `json:"model"`
	Stream    bool       `json:"stream"`
	Messages  []*message `json:"messages"`
	MaxTokens int        `json:"max_tokens"`
}

func newhHuggingfaceRequest(model string, messages []*message, maxTokens int) *huggingfaceRequest {
	return &huggingfaceRequest{
		Model:     model,
		Messages:  messages,
		MaxTokens: maxTokens,
		Stream:    false,
	}
}
