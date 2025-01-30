package ports

import "kororo/internal/core/domain/models"

type FieldDetectorService interface {
	DetectFields(intention *models.Intention, text string) ([]models.FieldValue, error)
}
