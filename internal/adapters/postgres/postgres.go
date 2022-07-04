package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"gitlab.com/g6834/team17/analytics-service/internal/config"
	"net"
)

type Database struct {
	pool *pgxpool.Pool
}

func New() *Database {
	return &Database{pool: nil}
}

func (d *Database) Connect(ctx context.Context) error {
	var err error

	conf := config.GetConfig()
	connstr := fmt.Sprintf("postgres://%s:%s@%s/%s",
		conf.DB.User,
		conf.DB.Password,
		net.JoinHostPort(conf.DB.Host, conf.DB.Port),
		conf.DB.DbName,
	)

	connConf, err := pgxpool.ParseConfig(connstr)
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
