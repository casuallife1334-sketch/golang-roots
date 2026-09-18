package repository

import (
	"context"
	"errors"
	"genealogy-tree/internal/core/domain"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/oklog/ulid/v2"
)

func (r *UsersRepository) CreateUser(ctx context.Context, input domain.CreateUserInput) (domain.User, error) {
	var user domain.User
	err := r.db.QueryRow(ctx, `
		INSERT INTO users (id, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, email, created_at, updated_at
	`, ulid.Make().String(), input.Email, input.PasswordHash).Scan(
		&user.ID, &user.Email, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.User{}, ErrDuplicate
		}
		return domain.User{}, err
	}
	return user, nil
}
