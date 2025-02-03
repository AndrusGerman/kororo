package ports

import (
	"context"
	"kororo/internal/core/domain/models"
)

// IntentionRepository ...
type IntentionRepository interface {
	BaseRepository[models.Intention]
}

type IntentionService interface {
	Detect(ctx context.Context, text string) (*models.Intention, error)
}

type IntentionProccessService interface {
	Process(ctx context.Context, text string) (string, error)
}

type MultiIntentionProccessService interface {
	Process(ctx context.Context, text string) (string, error)
}
