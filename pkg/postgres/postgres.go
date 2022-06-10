package postgres

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
	"os"
)

func New() *pgxpool.Pool {
	conConf, err := pgxpool.ParseConfig(os.Getenv("PG_CONNSTR"))
	if err != nil {
		log.Fatalf("unable parse config from ENV")
	}
	conn, err := pgxpool.ConnectConfig(context.TODO(), conConf)
	return conn
}
