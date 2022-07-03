package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"os"
)

type Database struct {
	pool *pgxpool.Pool
}

func New() *Database {
	return &Database{pool: nil}
}

func (d *Database) Connect(ctx context.Context) error {
	var err error

	connConf, err := pgxpool.ParseConfig(os.Getenv("PG_CONNSTR"))
	if err != nil {
		return fmt.Errorf("error parsing Postgres connstr: %w", err)
	}

	d.pool, err = pgxpool.ConnectConfig(ctx, connConf)
	if err != nil {
		return fmt.Errorf("error creating connections Pool: %w", err)
	}

	return nil
}

func (d *Database) Close(ctx context.Context) error {
	if d.pool != nil {
		d.pool.Close()
	}

	return nil
}
