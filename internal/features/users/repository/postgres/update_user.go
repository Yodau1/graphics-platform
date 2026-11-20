package users_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Yodau1/graphics-platform/internal/core/domain"
	core_errors "github.com/Yodau1/graphics-platform/internal/core/errors"
	core_postgres_pool "github.com/Yodau1/graphics-platform/internal/core/repository/postgres/pool"
)

func (r *UsersRepository) UpdateUser(ctx context.Context, user domain.User) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		UPDATE graphics.users
		SET name = $1,
			email = $2,
			version = version + 1
		WHERE id = $3 and version = $4
		RETURNING id, name, email, role, is_active, version;
	`

	row := r.pool.QueryRow(
		ctx,
		query,
		user.Name,
		user.Email,
		user.ID,
		user.Version,
	)

	var model userModel
	if err := model.Scan(row); err != nil {
		switch {
		case errors.Is(err, core_postgres_pool.ErrNoRows):
			return domain.User{}, fmt.Errorf(
				"user with id='%s' concurrently accessed: %w",
				user.ID,
				core_errors.ErrConflict,
			)
		case errors.Is(err, core_postgres_pool.ErrViolatesUnique):
			return domain.User{}, fmt.Errorf("email already registered: %w", core_errors.ErrConflict)
		default:
			return domain.User{}, fmt.Errorf("scan updated user: %w", err)
		}
	}

	userDomain := modelToDomain(model)

	return userDomain, nil
}
