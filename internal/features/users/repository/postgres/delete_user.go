package users_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	core_errors "github.com/Yodau1/graphics-platform/internal/core/errors"
	core_postgres_pool "github.com/Yodau1/graphics-platform/internal/core/repository/postgres/pool"
	"github.com/google/uuid"
)

func (r *UsersRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		DELETE FROM graphics.users
		WHERE id = $1;
	`

	commandTag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrViolatesForeignKey) {
			return fmt.Errorf("user has related records: %w", core_errors.ErrConflict)
		}
		return fmt.Errorf("delete user query: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("user with id='%s': %w", id, core_errors.ErrNotFound)
	}

	return nil
}
