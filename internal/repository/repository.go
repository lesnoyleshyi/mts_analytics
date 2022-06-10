package repository

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"mts_analytics/internal/domain"
	"mts_analytics/pkg/postgres"
)

type repo struct {
	Pool *pgxpool.Pool
}

func New() *repo {
	pool := postgres.New()
	return &repo{Pool: pool}
}

func (r repo) Save(event domain.Event) error {
	return nil
}

func (r repo) GetSignedCount() (int, error) {
	return 0, nil
}

func (r repo) GetNotSignedYetCount() (int, error) {
	return 0, nil
}

func (r repo) GetSignitionTotalTime(taskUUID string) (int, error) {
	return 0, nil
}
