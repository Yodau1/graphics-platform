package users_postgres_repository

import (
	"context"
	"fmt"

	"github.com/Yodau1/graphics-platform/internal/core/domain"
)

func (r *UsersRepository) GetUsers(
	ctx context.Context,
	limit int,
	offset int,
) ([]domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT id, name, email, role, is_active, version
		FROM graphics.users
		ORDER BY id ASC
		LIMIT $1 OFFSET $2;
	`

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("select users: %w", err)
	}
	defer rows.Close()

	var users []domain.User

	for rows.Next() {
		var model userModel
		if err := model.Scan(rows); err != nil {
			return nil, fmt.Errorf("scan user: %w", err)
		}

		users = append(users, modelToDomain(model))
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate users: %w", err)
	}

	return users, nil
}
