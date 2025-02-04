package huggingface

import (
	"errors"
	"fmt"
	"strings"

	"kororo/internal/adapters/config"
	"kororo/internal/core/domain/models"
	"kororo/internal/core/ports"
)

type huggingface struct {
	rest  ports.RestAdapter
	conf  *config.Config
	model string
}

func New(rest ports.RestAdapter, conf *config.Config, model string) ports.LLMAdapter {
	var llm = new(huggingface)
	llm.rest = rest
	llm.model = model
	return llm
}

func (o *huggingface) newMessages(base []*models.Message) []*message {
	var messages = make([]*message, len(base))

	for i := range base {
		if base[i].RoleID == models.AssistantRoleID {
			messages[i] = newMessage("assistant", base[i].Content)
		}

		if base[i].RoleID == models.UserRoleID {
			messages[i] = newMessage("user", base[i].Content)
		}

		if base[i].RoleID == models.SystemRoleID {
			messages[i] = newMessage("assistant", base[i].Content)
		}
	}
	return messages
}

func (o *huggingface) Quest(base []*models.Message) (*models.Message, error) {
	var messages = o.newMessages(base)
	var response *messageResponse

	var request = newhHuggingfaceRequest(o.model, messages, 2048)
	var err error

	if response, err = o.newRequest(request); err != nil {
		return nil, err
	}

	if len(response.Choices) == 0 {
		return nil, errors.New("no response")
	}

	return models.NewMessage(strings.TrimSpace(response.Choices[0].Message.Content), models.AssistantRoleID), nil
}

func (o *huggingface) ProcessSystemMessage(systemMessage string, userMessage string) (string, error) {

	var messages = []*models.Message{
		models.NewMessage(systemMessage, models.SystemRoleID),
		models.NewMessage(userMessage, models.UserRoleID),
	}

	var response, err = o.Quest(messages)
	if err != nil {
		return "", err
	}

	return response.Content, nil

}

func (o *huggingface) QuestParts(base []*models.Message, partsSize int) (<-chan *models.Message, error) {

	var messageStream = make(chan *models.Message)
	return messageStream, nil
}

func (o *huggingface) newRequest(request *huggingfaceRequest) (*messageResponse, error) {
	var response = new(messageResponse)

	var url = fmt.Sprintf("https://api-inference.huggingface.co/models/%s/v1/chat/completions", o.model)

	var err = o.rest.Post(url, map[string]string{
		"Authorization": "Bearer " + o.conf.HUGGINGFACE_API_KEY(),
		"Content-Type":  "application/json",
	}, request, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
