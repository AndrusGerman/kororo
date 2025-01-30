package models

type LLMError struct {
	UserMessage     string
	InternalMessage string
}

func (e *LLMError) Error() string {
	return e.InternalMessage
}

func (e *LLMError) GetUserMessage() string {
	return e.UserMessage
}

func NewLLMError(userMessage string, internalMessage string) error {
	return &LLMError{
		UserMessage:     userMessage,
		InternalMessage: internalMessage,
	}
}
