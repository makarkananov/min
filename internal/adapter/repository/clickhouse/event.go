package clickhouse

import (
	"context"
	"database/sql"
	"fmt"
	"min/internal/core/domain"
)

type EventRepository struct {
	db *sql.DB
}

func NewEventRepository(db *sql.DB) *EventRepository {
	return &EventRepository{db: db}
}

func (r *EventRepository) AddEvent(ctx context.Context, event domain.Event) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	stmt, err := tx.PrepareContext(
		ctx,
		`INSERT INTO events (short_url, original_url, timestamp, user_agent, ip) VALUES (?, ?, ?, ?, ?)`,
	)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to prepare insert statement: %w", err)
	}
	defer stmt.Close()

	if _, err := stmt.ExecContext(
		ctx,
		event.ShortURL,
		event.OriginalURL,
		event.Timestamp,
		event.UserAgent,
		event.IP,
	); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to execute insert statement: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
