package openai

import (
	"context"
	"strings"
	"sync"

	"kororo/internal/core/domain/models"
	"kororo/internal/core/ports"

	openai "github.com/openai/openai-go"
	openaioption "github.com/openai/openai-go/option"
)

type ollama struct {
	model  string
	mt     *sync.Mutex
	apiKey string
	url    string

	client *openai.Client
}

func New(model string, apiKey string, url string) ports.LLMAdapter {
	var llm = new(ollama)
	llm.mt = new(sync.Mutex)
	llm.model = model
	llm.apiKey = apiKey
	llm.url = url

	llm.client = openai.NewClient(
		openaioption.WithBaseURL(url),
		openaioption.WithAPIKey(apiKey),
	)
	return llm

}

func (o *ollama) ProcessSystemMessage(systemMessage string, userMessage string) (string, error) {
	var messages = make([]openai.ChatCompletionMessageParamUnion, 0)
	messages = append(messages, openai.SystemMessage(systemMessage))
	messages = append(messages, openai.UserMessage(userMessage))

	var request = newRequest(o.model, messages, false)

	chatCompletion, err := o.client.Chat.Completions.New(context.TODO(), request)
	if err != nil {
		return "", err
	}

	return chatCompletion.Choices[0].Message.Content, nil

}

func (o *ollama) newMessages(base []*models.Message) []openai.ChatCompletionMessageParamUnion {
	var messages = make([]openai.ChatCompletionMessageParamUnion, len(base))

	for i := range base {
		if base[i].RoleID == models.AssistantRoleID {
			messages[i] = openai.AssistantMessage(base[i].Content)
		}

		if base[i].RoleID == models.UserRoleID {
			messages[i] = openai.UserMessage(base[i].Content)
		}

		if base[i].RoleID == models.SystemRoleID {
			messages[i] = openai.SystemMessage(base[i].Content)
		}

	}
	return messages
}

func (o *ollama) Quest(base []*models.Message) (*models.Message, error) {
	var messages = o.newMessages(base)

	var request = newRequest(o.model, messages, false)

	chatCompletion, err := o.client.Chat.Completions.New(context.TODO(), request)
	if err != nil {
		return nil, err
	}

	return models.NewMessage(strings.TrimSpace(chatCompletion.Choices[0].Message.Content), models.AssistantRoleID), nil
}

func (o *ollama) QuestParts(base []*models.Message, partsSize int) (<-chan *models.Message, error) {
	//var messages = o.newMessages(base)
	//var response *messageResponse
	//var messageResponseStream <-chan *messageResponse
	var messageStream = make(chan *models.Message)
	var err error

	//var request = newOllamaRequest(o.model, messages, true)

	return messageStream, nil

	//if messageResponseStream, err = o.newRequestStream(nil); err != nil {
	//	return nil, err
	//}
	//
	//go func() {
	//	var text string
	//	for messageResponse := range messageResponseStream {
	//		text += messageResponse.Message.Content
	//
	//		if len(text) > partsSize || messageResponse.Done {
	//			messageStream <- models.NewMessage(strings.TrimSpace(text), models.AssistantRoleID)
	//			text = ""
	//		}
	//	}
	//
	//	close(messageStream)
	//}()

	return messageStream, err
}
