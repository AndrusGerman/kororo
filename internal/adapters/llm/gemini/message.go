package gemini

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
