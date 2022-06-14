package postgres

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
	"os"
)

const app_name = "analytics"
const host_ip = `127.0.0.1`

func New() *pgxpool.Pool {
	fn := "postgres_New"

	conConf, err := pgxpool.ParseConfig(os.Getenv("PG_CONNSTR"))
	if err != nil {
		log.WithFields(log.Fields{
			"app_name":    app_name,
			"host_ip":     host_ip,
			"logger_name": fn,
		}).Fatalf("unable parse config from ENV: %s", err)
	}
	conn, err := pgxpool.ConnectConfig(context.TODO(), conConf)
	if err != nil {
		log.WithFields(log.Fields{
			"app_name":    app_name,
			"host_ip":     host_ip,
			"logger_name": fn,
		}).Fatalf("uanble establish connection with database: %s", err)
	}
	return conn
}
