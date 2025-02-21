package services

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"kororo/internal/core/domain"
	"kororo/internal/core/domain/models"
	"kororo/internal/core/ports"
	"os"
	"strings"
)

type MultiIntentionChatProcessService struct {
	IntentionProccessService ports.IntentionProccessService
	llmAdapter               ports.LLMAdapter
	logger                   ports.LogService
}

func (m *MultiIntentionChatProcessService) Process(ctx context.Context, initialMessage string) error {
	var messages = make([]*models.Message, 0)
	// system prompt:
	messages = append(messages, models.NewMessageFromSystem(domain.MultiIntentionChatPrompt))

	// Agregar el mensaje del usuario
	messages = append(messages, m.NewMessageUserFromUser(initialMessage))

	for {
		var llmError = new(models.LLMError)

		// preguntar al sistema que quiere hacer
		systemResponse, err := m.Quest(ctx, messages)
		if err != nil {
			return err
		}

		messages = append(messages, models.NewMessageFromAssistant(systemResponse.Json()))

		// el sistema habla con el usuario
		if systemResponse.ToUser != "" {
			m.logger.Info("systemResponse.ToUser", systemResponse.ToUser)
			messages = append(messages, m.NewMessageUserFromUser(m.getUserInput()))
			continue
		}

		if systemResponse.ToSystem == "" {
			return errors.New("systemResponse.ToSystem is empty")
		}

		// el sistema quiere hacer una accion
		m.logger.Info("systemResponse.ToSystem", systemResponse.ToSystem)
		responseIntention, err := m.IntentionProccessService.Process(ctx, systemResponse.ToSystem)
		if err != nil {
			if errors.Is(err, domain.ErrIntentionNotFound) {
				m.logger.Info("IntentionProccessService.Err", "Esta accion no la puede realizar el sistema")
				messages = append(messages, m.NewMessageUserFromSystem("Esta accion no la puede realizar el sistema"))
				continue
			}

			if errors.Is(err, domain.ErrMultipleIntentionSend) {
				errors.As(err, &llmError)
				m.logger.Info("IntentionProccessService.Err", llmError.InternalMessage)
				messages = append(messages, m.NewMessageUserFromSystem(llmError.UserMessage))
				continue
			}

			return err
		}

		m.logger.Info("IntentionProccessService.ProcessResponse", responseIntention)
		messages = append(messages, m.NewMessageUserFromSystem(responseIntention))
	}

}

func (m *MultiIntentionChatProcessService) NewMessageUserFromSystem(message string) *models.Message {
	return models.NewMessageFromUser((&models.MultiIntentionInput{FromSystem: message}).Json())
}

func (m *MultiIntentionChatProcessService) NewMessageUserFromUser(message string) *models.Message {
	return models.NewMessageFromUser((&models.MultiIntentionInput{FromUser: message}).Json())
}

func (m *MultiIntentionChatProcessService) Quest(ctx context.Context, messages []*models.Message) (*models.MultiIntentionInput, error) {
	respAssistant, err := m.llmAdapter.Quest(messages)
	if err != nil {
		return nil, err
	}

	return models.NewMultiIntentionInputFromString(respAssistant.Content)
}

func (m *MultiIntentionChatProcessService) getUserInput() string {
	var reader = bufio.NewReader(os.Stdin)
	fmt.Print("Prompt: ")
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)
	return text
}

func NewMultiIntentionChatProcessService(intentionProccessService ports.IntentionProccessService, llmAdapter ports.LLMAdapter, logger ports.LogService) *MultiIntentionChatProcessService {
	return &MultiIntentionChatProcessService{
		IntentionProccessService: intentionProccessService,
		llmAdapter:               llmAdapter,
		logger:                   logger,
	}
}
