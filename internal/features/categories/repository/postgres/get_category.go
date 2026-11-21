package categories_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/Yodau1/graphics-platform/internal/core/domain"
	core_errors "github.com/Yodau1/graphics-platform/internal/core/errors"
	core_postgres_pool "github.com/Yodau1/graphics-platform/internal/core/repository/postgres/pool"
)

// GetCategory возвращает категорию по ID вместе с её текущей версией
func (r *CategoriesRepository) GetCategory(
	ctx context.Context,
	id uuid.UUID,
) (domain.Category, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT id, name, version
		FROM graphics.categories
		WHERE id = $1;
	`

	row := r.pool.QueryRow(ctx, query, id)

	var model CategoryModel
	if err := model.Scan(row); err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.Category{}, fmt.Errorf("category with id='%s': %w", id, core_errors.ErrNotFound)
		}

		return domain.Category{}, fmt.Errorf("scan category: %w", err)
	}

	return modelToDomain(model), nil
}
