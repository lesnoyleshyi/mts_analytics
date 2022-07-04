package app

import (
	"context"
	myMongo "gitlab.com/g6834/team17/analytics-service/internal/adapters/mongo"
	"gitlab.com/g6834/team17/analytics-service/internal/adapters/postgres"
	ports "gitlab.com/g6834/team17/analytics-service/internal/ports/output"
)

type Storage interface {
	Connect(ctx context.Context) error
	Close(ctx context.Context) error
	ports.EventStorage
}

func NewStorage(storageType string) Storage {
	switch storageType {
	case "mongo":
		return myMongo.New()
	case "postgres":
		return postgres.New()
	default:
		return postgres.New()
	}
}
