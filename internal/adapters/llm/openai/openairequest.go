package openai

import (
	openai "github.com/openai/openai-go"
)

func newRequest(model string, messages []openai.ChatCompletionMessageParamUnion) openai.ChatCompletionNewParams {

	return openai.ChatCompletionNewParams{
		Model:     openai.F(model),
		Messages:  openai.F(messages),
		MaxTokens: openai.Int(1500),
	}
}
