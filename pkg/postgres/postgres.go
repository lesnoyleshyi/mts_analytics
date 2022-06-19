package postgres

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
	"os"
)

const appName = "analytics"
const hostIP = `127.0.0.1`

func New() *pgxpool.Pool {
	fn := "postgres_New"

	conConf, err := pgxpool.ParseConfig(os.Getenv("PG_CONNSTR"))
	if err != nil {
		log.WithFields(log.Fields{
			"appName":     appName,
			"hostIP":      hostIP,
			"logger_name": fn,
		}).Fatalf("unable parse config from ENV: %s", err)
	}

	conn, err := pgxpool.ConnectConfig(context.TODO(), conConf)
	if err != nil {
		log.WithFields(log.Fields{
			"appName":     appName,
			"hostIP":      hostIP,
			"logger_name": fn,
		}).Fatalf("uanble establish connection with database: %s", err)
	}

	return conn
}
