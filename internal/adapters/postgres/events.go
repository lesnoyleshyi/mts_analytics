package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"gitlab.com/g6834/team17/analytics-service/internal/domain/entity"
)

const errBeginTx = `error begin transaction`
const errCommitTx = `error commit transaction`
const errAddToDB = `error adding data to db`
const errRetrieveFromDB = `error retrieving data from db`

const saveQuery = `INSERT INTO task_events (task_uuid, event, user_uuid, timestamp)
					VALUES ($1, $2, $3, $4);`

func (d Database) Save(ctx context.Context, e entity.Event) error {
	// Probably we should return error in this case
	if e.UserUUID == "" {
		e.UserUUID = "00000000-0000-0000-0000-000000000000"
	}

	txOpts := pgx.TxOptions{
		IsoLevel:       pgx.ReadUncommitted,
		AccessMode:     pgx.ReadWrite,
		DeferrableMode: pgx.NotDeferrable,
	}
	tx, err := d.Pool.BeginTx(ctx, txOpts)
	defer func() { _ = tx.Rollback(ctx) }()
	if err != nil {
		return fmt.Errorf("%s: %w", errBeginTx, err)
	}

	res, err := tx.Exec(ctx, saveQuery, e.TaskUUID, e.EventType, e.UserUUID, e.Timestamp)
	if err != nil || res.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", errAddToDB, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("%s: %w", errCommitTx, err)
	}

	return nil
}

const getSignedCountQuery = `SELECT count(task_uuid) FROM task_events WHERE
							event = 'signed';`

func (d Database) GetSignedCount(ctx context.Context) (uint, error) {
	var count uint

	txOpts := pgx.TxOptions{
		IsoLevel:       pgx.ReadCommitted,
		AccessMode:     pgx.ReadOnly,
		DeferrableMode: pgx.NotDeferrable,
	}

	tx, err := d.Pool.BeginTx(ctx, txOpts)
	defer func() { _ = tx.Rollback(ctx) }()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", errBeginTx, err)
	}

	row := tx.QueryRow(ctx, getSignedCountQuery)
	if err := row.Scan(&count); err != nil {
		return 0, fmt.Errorf("%s: %w", errRetrieveFromDB, err)
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

func (d Database) GetUnsignedCount(ctx context.Context) (uint, error) {
	var count uint

	txOpts := pgx.TxOptions{
		IsoLevel:       pgx.ReadCommitted,
		AccessMode:     pgx.ReadOnly,
		DeferrableMode: pgx.NotDeferrable,
	}

	tx, err := d.Pool.BeginTx(ctx, txOpts)
	defer func() { _ = tx.Rollback(ctx) }()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", errBeginTx, err)
	}

	row := tx.QueryRow(ctx, getUnsignedCountQuery)
	if err := row.Scan(&count); err != nil {
		return 0, fmt.Errorf("%s: %w", errRetrieveFromDB, err)
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

func (d Database) GetSignitionTime(ctx context.Context, event entity.Event) (uint64, error) {
	var Sec float64

	txOpts := pgx.TxOptions{
		IsoLevel:       pgx.ReadCommitted,
		AccessMode:     pgx.ReadOnly,
		DeferrableMode: pgx.NotDeferrable,
	}
	tx, err := d.Pool.BeginTx(ctx, txOpts)
	defer func() { _ = tx.Rollback(ctx) }()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", errBeginTx, err)
	}

	row := tx.QueryRow(ctx, getSignTimeQuery, event.TaskUUID)
	if err := row.Scan(&Sec); err != nil {
		return 0, fmt.Errorf("%s: %w", errRetrieveFromDB, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("%s: %w", errCommitTx, err)
	}

	return uint64(Sec), nil
}
