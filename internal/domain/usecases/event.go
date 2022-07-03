package usecases

import (
	"context"
	"gitlab.com/g6834/team17/analytics-service/internal/domain/entity"
	ports "gitlab.com/g6834/team17/analytics-service/internal/ports/output"
)

type EventService struct {
	storage ports.EventStorage
}

func NewEventService(storage ports.EventStorage) EventService {
	return EventService{storage: storage}
}

func (e EventService) Save(ctx context.Context, event entity.Event) error {
	return e.storage.Save(ctx, event)
}

func (e EventService) GetSignedCount(ctx context.Context) (uint, error) {
	return e.storage.GetSignedCount(ctx)
}

func (e EventService) GetUnsignedCount(ctx context.Context) (uint, error) {
	return e.storage.GetUnsignedCount(ctx)
}

func (e EventService) GetSignitionTime(ctx context.Context, event entity.Event) (uint64, error) {
	return e.storage.GetSignitionTime(ctx, event)
}
