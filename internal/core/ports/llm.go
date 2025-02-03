package ports

import "kororo/internal/core/domain/models"

type LLMAdapter interface {
	Quest(base []*models.Message) (*models.Message, error)
	QuestParts(base []*models.Message, partsSize int) (<-chan *models.Message, error)
	ProcessSystemMessage(systemMessage string, text string) (string, error)
}
