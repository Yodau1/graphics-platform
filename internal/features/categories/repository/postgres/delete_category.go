package categories_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	core_errors "github.com/Yodau1/graphics-platform/internal/core/errors"
	core_postgres_pool "github.com/Yodau1/graphics-platform/internal/core/repository/postgres/pool"
)

func (r *CategoriesRepository) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `DELETE FROM graphics.categories WHERE id = $1;`

	tag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrViolatesForeignKey) {
			return fmt.Errorf("category is used by resources: %w", core_errors.ErrConflict)
		}

		return fmt.Errorf("delete category query: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("category with id='%s': %w", id, core_errors.ErrNotFound)
	}

	return nil
}
