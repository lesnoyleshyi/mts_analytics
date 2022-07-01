package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Database struct {
	Pool *pgxpool.Pool
}

func New(ctx context.Context, connstr string) (Database, error) {
	connConf, err := pgxpool.ParseConfig(connstr)
	if err != nil {
		return Database{}, fmt.Errorf("error parsing Postgres connstr: %w", err)
	}

	pool, err := pgxpool.ConnectConfig(ctx, connConf)
	if err != nil {
		return Database{}, fmt.Errorf("error creating connections Pool: %w", err)
	}

	return Database{Pool: pool}, nil
}
