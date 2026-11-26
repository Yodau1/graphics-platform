package users_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/Yodau1/graphics-platform/internal/core/domain"
	core_errors "github.com/Yodau1/graphics-platform/internal/core/errors"
	core_postgres_pool "github.com/Yodau1/graphics-platform/internal/core/repository/postgres/pool"
)

func (r *UsersRepository) UpdateUserStatus(
	ctx context.Context,
	id uuid.UUID,
	version int64,
	isActive bool,
) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		UPDATE graphics.users
		SET is_active = $1,
		    version = version + 1
		WHERE id = $2 AND version = $3
		RETURNING id, name, email, role, is_active, version;
	`

	row := r.pool.QueryRow(ctx, query, isActive, id, version)

	var model userModel
	if err := model.Scan(row); err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.User{}, fmt.Errorf(
				"user with id='%s' concurrently accessed: %w",
				id,
				core_errors.ErrConflict,
			)
		}
		return domain.User{}, fmt.Errorf("scan updated user status: %w", err)
	}

	return modelToDomain(model), nil
}
