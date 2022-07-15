package mongo

import (
	"context"
	"fmt"
	"gitlab.com/g6834/team17/analytics-service/internal/domain/entity"
)

func (d Database) Save(ctx context.Context, e entity.Event) error {
	return fmt.Errorf("implement me") //nolint:goerr113
}

func (d Database) GetSignedCount(ctx context.Context) (uint, error) {
	return 0, fmt.Errorf("implement me") // nolint:goerr113
}

func (d Database) GetUnsignedCount(ctx context.Context) (uint, error) {
	return 0, fmt.Errorf("implement me") // nolint:goerr113
}

func (d Database) GetSignitionTime(ctx context.Context, event entity.Event) (uint64, error) {
	return 0, fmt.Errorf("implement me") // nolint:goerr113
}
