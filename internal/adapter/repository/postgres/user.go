package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"min/internal/core/domain"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Save(ctx context.Context, user *domain.User) error {
	_, err := r.db.ExecContext(
		ctx,
		"INSERT INTO users (username, password, role, plan_name, links_remaining) VALUES ($1, $2, $3, $4, $5)",
		user.Username,
		user.Password,
		user.Role,
		user.Plan,
		user.LinksRemaining,
	)

	if err != nil {
		return fmt.Errorf("error saving user: %w", err)
	}

	return nil
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	row := r.db.QueryRowContext(
		ctx,
		"SELECT username, password, role, plan_name, links_remaining FROM users WHERE username = $1",
		username,
	)
	var user domain.User
	err := row.Scan(&user.Username, &user.Password, &user.Role, &user.Plan, &user.LinksRemaining)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("error getting user by username: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) ChangeLinksRemaining(ctx context.Context, username string, linksRemaining int64) error {
	_, err := r.db.ExecContext(
		ctx,
		"UPDATE users SET links_remaining = $1 WHERE username = $2",
		linksRemaining,
		username,
	)

	if err != nil {
		return fmt.Errorf("error changing links remaining: %w", err)
	}

	return nil
}
