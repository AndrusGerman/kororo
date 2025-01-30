package repository

import (
	mongodb "kororo/internal/adapters/storage/mongo"
	"kororo/internal/core/domain/models"
	"kororo/internal/core/ports"
)

func NewIntentionRepository(mongoService *mongodb.Mongo) ports.IntentionRepository {
	return &IntentionRepository{
		BaseRepository: newBaseRepository[models.Intention](mongoService, "intentions"),
	}
}

type IntentionRepository struct {
	ports.BaseRepository[models.Intention]
}
