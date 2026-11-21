package categories_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Yodau1/graphics-platform/internal/core/domain"
	core_errors "github.com/Yodau1/graphics-platform/internal/core/errors"
	core_postgres_pool "github.com/Yodau1/graphics-platform/internal/core/repository/postgres/pool"
)

func (r *CategoriesRepository) UpdateCategory(
	ctx context.Context,
	category domain.Category,
) (domain.Category, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		UPDATE graphics.categories
		SET name = $1,
		    version = version + 1
		WHERE id = $2 AND version = $3
		RETURNING id, name, version;
	`

	row := r.pool.QueryRow(ctx, query, category.Name, category.ID, category.Version)

	var model CategoryModel
	if err := model.Scan(row); err != nil {
		switch {
		case errors.Is(err, core_postgres_pool.ErrNoRows):
			// После чтения запись могла измениться или быть удалена.
			return domain.Category{}, fmt.Errorf(
				"category with id='%s' changed concurrently: %w",
				category.ID,
				core_errors.ErrConflict,
			)
		case errors.Is(err, core_postgres_pool.ErrViolatesUnique):
			return domain.Category{}, fmt.Errorf("category name already exists: %w", core_errors.ErrConflict)
		default:
			return domain.Category{}, fmt.Errorf("scan updated category: %w", err)
		}
	}

	return modelToDomain(model), nil
}
