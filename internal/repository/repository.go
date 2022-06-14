package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
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

const errBeginTx = `error begin transaction`
const errCommitTx = `error commit transaction`
const errAddToDb = `error adding data to db`
const errRetrieveFromDb = `error retrieving data from db`

const saveQuery = `INSERT INTO task_events (task_uuid, event, user_uuid, timestamp)
					VALUES ($1, $2, $3, $4);`

func (r repo) Save(e domain.Event) error {
	if e.UserUUID == "" {
		e.UserUUID = "00000000-0000-0000-0000-000000000000"
	}
	ctx := context.TODO()
	txOpts := pgx.TxOptions{
		IsoLevel:       pgx.ReadUncommitted,
		AccessMode:     pgx.ReadWrite,
		DeferrableMode: pgx.NotDeferrable,
	}
	tx, err := r.Pool.BeginTx(ctx, txOpts)
	defer func() { _ = tx.Rollback(ctx) }()
	if err != nil {
		return fmt.Errorf("%s: %w", errBeginTx, err)
	}

	res, err := tx.Exec(ctx, saveQuery, e.TaskUUID, e.EventType, e.UserUUID, e.Timestamp)
	if err != nil || res.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", errAddToDb, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("%s: %w", errCommitTx, err)
	}

	return nil
}

const getSignedCountQuery = `SELECT count(task_uuid) FROM task_events WHERE
							event = 'signed';`

func (r repo) GetSignedCount() (int, error) {
	ctx := context.TODO()
	var count int
	txOpts := pgx.TxOptions{
		IsoLevel:       pgx.ReadCommitted,
		AccessMode:     pgx.ReadOnly,
		DeferrableMode: pgx.NotDeferrable,
	}
	tx, err := r.Pool.BeginTx(ctx, txOpts)
	defer func() { _ = tx.Rollback(ctx) }()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", errBeginTx, err)
	}

	row := tx.QueryRow(ctx, getSignedCountQuery)
	if err := row.Scan(&count); err != nil {
		return 0, fmt.Errorf("%s: %w", errRetrieveFromDb, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("%s: %w", errCommitTx, err)
	}

	return count, nil
}

const getUnsignedCountQuery = `(SELECT (total - signed) FROM
								(SELECT
								count(task_uuid) FILTER (WHERE event = 'created')
										as total,
								count(task_uuid) FILTER (WHERE event = 'signed')
										as signed
								FROM task_events
								) AS unsigned);`

func (r repo) GetNotSignedYetCount() (int, error) {
	ctx := context.TODO()
	var count int
	txOpts := pgx.TxOptions{
		IsoLevel:       pgx.ReadCommitted,
		AccessMode:     pgx.ReadOnly,
		DeferrableMode: pgx.NotDeferrable,
	}
	tx, err := r.Pool.BeginTx(ctx, txOpts)
	defer func() { _ = tx.Rollback(ctx) }()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", errBeginTx, err)
	}

	row := tx.QueryRow(ctx, getUnsignedCountQuery)
	if err := row.Scan(&count); err != nil {
		return 0, fmt.Errorf("%s: %w", errRetrieveFromDb, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("%s: %w", errCommitTx, err)
	}

	return count, nil
}

const getSignTimeQuery = `SELECT EXTRACT(EPOCH FROM
			(SELECT (max - min) FROM
			(SELECT max(timestamp) AS max, min(timestamp) AS min FROM
				(SELECT timestamp FROM task_events WHERE task_uuid = $1) AS task_ts
			) AS total_time));`

func (r repo) GetSignitionTotalTime(taskUUID string) (int, error) {
	ctx := context.TODO()
	var Sec float64
	txOpts := pgx.TxOptions{
		IsoLevel:       pgx.ReadCommitted,
		AccessMode:     pgx.ReadOnly,
		DeferrableMode: pgx.NotDeferrable,
	}
	tx, err := r.Pool.BeginTx(ctx, txOpts)
	defer func() { _ = tx.Rollback(ctx) }()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", errBeginTx, err)
	}

	row := tx.QueryRow(ctx, getSignTimeQuery, taskUUID)
	if err := row.Scan(&Sec); err != nil {
		return 0, fmt.Errorf("%s: %w", errRetrieveFromDb, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("%s: %w", errCommitTx, err)
	}

	return int(Sec), nil
}
