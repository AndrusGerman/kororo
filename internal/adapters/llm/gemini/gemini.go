package gemini

import (
	"context"
	"strings"

	"kororo/internal/adapters/config"
	"kororo/internal/core/domain/models"
	"kororo/internal/core/ports"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type gemini struct {
	conf   *config.Config
	client *genai.Client
	model  string
}

func New(conf *config.Config) (ports.LLMAdapter, error) {
	var llm = new(gemini)
	var err error
	llm.conf = conf

	llm.model = "gemini-1.5-flash"
	llm.model = "gemini-2.0-flash-exp"

	llm.client, err = genai.NewClient(context.Background(), option.WithAPIKey(conf.GEMINI_API_KEY()))

	return llm, err
}

func (o *gemini) BasicQuest(text string) (string, error) {
	return "", nil
}

func (o *gemini) ProcessSystemMessage(systemMessage string, userMessage string) (string, error) {
	model := o.client.GenerativeModel(o.model)
	return newGeminiSystemRequest(model, systemMessage, userMessage)
}

func (o *gemini) newMessages(base []*models.Message) []*message {
	var messages = make([]*message, len(base))

	for i := range base {
		if base[i].RoleID == models.AssistantRoleID {
			messages[i] = newMessage("model", base[i].Content)
		}

		if base[i].RoleID == models.UserRoleID {
			messages[i] = newMessage("user", base[i].Content)
		}

		if base[i].RoleID == models.SystemRoleID {
			messages[i] = newMessage("model", base[i].Content)
			//messages[i] = newMessage("system", base[i].Content)

		}
	}
	return messages
}

func (o *gemini) Quest(base []*models.Message) (*models.Message, error) {
	var messages = o.newMessages(base)

	model := o.client.GenerativeModel(o.model)

	var response, err = newGeminiRequest(model, messages)

	return models.NewMessage(strings.TrimSpace(response), models.AssistantRoleID), err
}

func (o *gemini) QuestParts(base []*models.Message, partsSize int) (<-chan *models.Message, error) {
	var messages = o.newMessages(base)
	//var response *messageResponse
	//var messageResponseStream <-chan *messageResponse
	var messageStream = make(chan *models.Message)

	model := o.client.GenerativeModel(o.model)
	var response, err = newGeminiRequest(model, messages)
	if err != nil {
		return nil, err
	}

	go func() {

		messageStream <- models.NewMessage(response, models.AssistantRoleID)
		close(messageStream)
	}()

	return messageStream, err
}
