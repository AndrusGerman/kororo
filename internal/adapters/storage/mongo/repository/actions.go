package repository

import (
	mongodb "kororo/internal/adapters/storage/mongo"
	"kororo/internal/core/domain/models"
	"kororo/internal/core/ports"
)

func NewActionRepository(mongoService *mongodb.Mongo) ports.ActionRepository {
	return &ActionRepository{
		BaseRepository: newBaseRepository[models.Action](mongoService, "actions"),
	}
}

type ActionRepository struct {
	ports.BaseRepository[models.Action]
}
