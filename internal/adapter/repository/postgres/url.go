package postgres

import (
	"context"
	"database/sql"
	"errors"
)

type URLRepository struct {
	db *sql.DB
}

func NewURLRepository(db *sql.DB) *URLRepository {
	return &URLRepository{
		db: db,
	}
}

func (r *URLRepository) GetOriginal(ctx context.Context, short string) (string, error) {
	var original string
	err := r.db.QueryRowContext(ctx, "SELECT original_url FROM url WHERE short_url = $1", short).Scan(&original)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}

		return "", err
	}

	return original, nil
}

func (r *URLRepository) Add(ctx context.Context, short, original string) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO url (short_url, original_url) VALUES ($1, $2)", short, original)
	if err != nil {
		return err
	}

	return nil
}

func (r *URLRepository) Remove(ctx context.Context, short string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM url WHERE short_url = $1", short)
	if err != nil {
		return err
	}

	return nil
}
