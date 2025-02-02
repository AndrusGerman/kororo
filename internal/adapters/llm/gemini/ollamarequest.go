package gemini

import (
	"context"
	"fmt"

	"github.com/google/generative-ai-go/genai"
)

func newGeminiRequest(model *genai.GenerativeModel, messages []*message) (string, error) {
	cs := model.StartChat()
	cs.History = []*genai.Content{}

	for ind, message := range messages {
		if ind == len(messages)-1 {
			break
		}
		cs.History = append(cs.History, &genai.Content{
			Parts: []genai.Part{
				genai.Text(message.Content),
			},
			Role: message.Role,
		})
	}

	resp, err := cs.SendMessage(context.TODO(), genai.Text(messages[len(messages)-1].Content))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s", resp.Candidates[0].Content.Parts[0]), nil

}

func newGeminiSystemRequest(model *genai.GenerativeModel, systemMessage string, userMessage string) (string, error) {
	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{
			genai.Text(systemMessage),
		},
	}

	var resp, err = model.GenerateContent(context.TODO(), genai.Text(userMessage))
	if err != nil {
		return "nil", err
	}

	return fmt.Sprintf("%s", resp.Candidates[0].Content.Parts[0]), nil

}
