package services

import (
	"context"
	"kororo/internal/core/domain"
	"kororo/internal/core/domain/models"
	"kororo/internal/core/ports"
)

type MultiIntentionProcessService struct {
	IntentionProccessService ports.IntentionProccessService
	llmAdapter               ports.LLMAdapter
}

func (m *MultiIntentionProcessService) Process(ctx context.Context, text string) (string, error) {
	var messages = make([]*models.Message, 0)

	// primera pregunta, conversacion systema a systema
	messages = append(messages, models.NewMessageFromSystem(domain.MultiIntentionPrompt))
	messages = append(messages, models.NewMessageFromUser((&models.MultiIntentionInput{FromUser: text}).Json()))

	resp, err := m.llmAdapter.Quest(messages)
	if err != nil {
		return "", err
	}

	messages = append(messages, models.NewMessageFromAssistant(resp.Content))

	// respuesta del sistema
	firstMessage, err := models.NewMultiIntentionInputFromString(resp.Content)
	if err != nil {
		return "", err
	}

	responseIntention, err := m.IntentionProccessService.Process(ctx, firstMessage.ToSystem)
	if err != nil {
		return "", err
	}

	messages = append(messages, models.NewMessageFromUser((&models.MultiIntentionInput{FromSystem: responseIntention}).Json()))

	for {

		resp, err := m.llmAdapter.Quest(messages)
		if err != nil {
			return "", err
		}

		systemResponse, err := models.NewMultiIntentionInputFromString(resp.Content)
		if err != nil {
			return "", err
		}

		if systemResponse.Finish {
			return systemResponse.ToUser, nil
		}

		messages = append(messages, models.NewMessageFromAssistant(systemResponse.Json()))
		responseIntention, err := m.IntentionProccessService.Process(ctx, systemResponse.ToSystem)
		if err != nil {
			return "", err
		}

		messages = append(messages, models.NewMessageFromUser((&models.MultiIntentionInput{FromSystem: responseIntention}).Json()))

	}

}

func NewMultiIntentionProcessService(intentionProccessService ports.IntentionProccessService, llmAdapter ports.LLMAdapter) ports.MultiIntentionProccessService {
	return &MultiIntentionProcessService{
		IntentionProccessService: intentionProccessService,
		llmAdapter:               llmAdapter,
	}
}
