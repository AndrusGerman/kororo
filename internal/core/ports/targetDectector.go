package ports

import "context"

type TargetDectector interface {
	Detect(ctx context.Context, text string) (string, error)
}
