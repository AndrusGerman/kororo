package ports

import (
	"context"
	"kororo/internal/core/domain/models"
	"kororo/internal/core/domain/types"
)

// ActionRepository ...
type ActionRepository interface {
	BaseRepository[models.Action]
}

type ActionService interface {
	GetAction(ctx context.Context, id types.Id) (*models.Action, error)
	ProcessAction(ctx context.Context, action *models.Action, actionContext *models.ActionPipelineContext) (*models.ActionResponse, error)
}
