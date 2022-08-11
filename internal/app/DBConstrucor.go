package app

import (
	"gitlab.com/g6834/team17/analytics-service/internal/adapters/interfaces"
	"gitlab.com/g6834/team17/analytics-service/internal/adapters/postgres"
)

func NewStorage(storageType string) interfaces.Storage { //nolint:ireturn
	switch storageType {
	case "postgres":
		return postgres.New()
	}

	return nil
}
