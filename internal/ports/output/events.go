package ports

import (
	"context"
	"gitlab.com/g6834/team17/analytics-service/internal/domain/entity"
)

type EventStorage interface {
	Save(ctx context.Context, event entity.Event) error
	GetSignedCount(ctx context.Context) (uint, error)
	GetUnsignedCount(ctx context.Context) (uint, error)
	GetSignitionTime(ctx context.Context, event entity.Event) (uint64, error)
}
