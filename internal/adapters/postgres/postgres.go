package postgres

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pressly/goose/v3"
	"log"
	"os"
	"gitlab.com/g6834/team17/analytics-service/internal/config"
	"net"
)

type Database struct {
	pool *pgxpool.Pool
}

func New() *Database {
	return &Database{pool: nil}
}

//go:embed changelog/*.sql
var embedMigrations embed.FS

func (d *Database) Connect(ctx context.Context) error {
	var err error

	connstr := os.Getenv("PG_CONNSTR")

	connConf, err := pgxpool.ParseConfig(connstr)
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

	if err = applyMigrations(connstr); err != nil {
		return fmt.Errorf("unable apply migrations: %w", err)
	}

	return nil
}

func (d *Database) Close(ctx context.Context) error {
	if d.pool != nil {
		d.pool.Close()
	}

	return nil
}

func applyMigrations(connstr string) error {
	goose.SetBaseFS(embedMigrations)

	dbConn, err := sql.Open("pgx", connstr)
	if err != nil {
		return err
	}
	defer func() {
		if err := dbConn.Close(); err != nil {
			log.Println(err)
		}
	}()

	if err = goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err = goose.Up(dbConn, "changelog"); err != nil {
		return err
	}

	return nil
}
