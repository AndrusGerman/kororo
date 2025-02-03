package services

import (
	"context"
	"errors"
	"kororo/internal/core/domain"
	"kororo/internal/core/domain/models"
	"kororo/internal/core/ports"
)

type MultiIntentionProcessService struct {
	IntentionProccessService ports.IntentionProccessService
	llmAdapter               ports.LLMAdapter
	logger                   ports.LogService
}

func (m *MultiIntentionProcessService) Process(ctx context.Context, text string) (string, error) {
	var messages = make([]*models.Message, 0)

	// system prompt:
	var systemMessage = models.NewMessageFromSystem(domain.MultiIntentionPrompt)
	messages = append(messages, systemMessage)

	// Agregar el mensaje del usuario
	var userMessage = models.NewMessageFromUser((&models.MultiIntentionInput{FromUser: text}).Json())
	messages = append(messages, userMessage)

	for {
		// preguntar al sistema que quiere hacer
		systemResponse, err := m.Quest(ctx, messages)
		if err != nil {
			return "", err
		}

		messages = append(messages, models.NewMessageFromAssistant(systemResponse.Json()))

		// el sistema quiere terminar
		if systemResponse.Finish || systemResponse.ToUser != "" {
			return systemResponse.ToUser, nil
		}

		if systemResponse.ToSystem == "" {
			return "", errors.New("systemResponse.ToSystem is empty")
		}

		// el sistema quiere hacer una accion
		m.logger.Info("systemResponse.ToSystem", systemResponse.ToSystem)
		responseIntention, err := m.IntentionProccessService.Process(ctx, systemResponse.ToSystem)
		if err != nil {
			return "", err
		}

		m.logger.Info("IntentionProccessService.ProcessResponse", responseIntention)
		messages = append(messages, models.NewMessageFromUser((&models.MultiIntentionInput{FromSystem: responseIntention}).Json()))

	}

}

func (m *MultiIntentionProcessService) Quest(ctx context.Context, messages []*models.Message) (*models.MultiIntentionInput, error) {
	respAssistant, err := m.llmAdapter.Quest(messages)
	if err != nil {
		return nil, err
	}

	return models.NewMultiIntentionInputFromString(respAssistant.Content)
}

func NewMultiIntentionProcessService(intentionProccessService ports.IntentionProccessService, llmAdapter ports.LLMAdapter, logger ports.LogService) ports.MultiIntentionProccessService {
	return &MultiIntentionProcessService{
		IntentionProccessService: intentionProccessService,
		llmAdapter:               llmAdapter,
		logger:                   logger,
	}
}
