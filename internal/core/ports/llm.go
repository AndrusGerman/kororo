package ports

import "kororo/internal/core/domain/models"

type LLMAdapter interface {
	BasicQuest(text string) (string, error)
	Quest(base []*models.Message, text string) (*models.Message, error)
	QuestParts(base []*models.Message, text string, partsSize int) (<-chan *models.Message, error)
	ProcessSystemMessage(systemMessage string, text string) (string, error)
}
