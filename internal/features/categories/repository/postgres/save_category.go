package categories_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Yodau1/graphics-platform/internal/core/domain"
	core_errors "github.com/Yodau1/graphics-platform/internal/core/errors"
	core_postgres_pool "github.com/Yodau1/graphics-platform/internal/core/repository/postgres/pool"
)

// SaveCategory вставляет новую категорию и возвращает сохранённую запись.
func (r *CategoriesRepository) SaveCategory(
	ctx context.Context,
	category domain.Category,
) (domain.Category, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		INSERT INTO graphics.categories (id, name, version)
		VALUES ($1, $2, $3)
		RETURNING id, name, version;
	`

	row := r.pool.QueryRow(ctx, query, category.ID, category.Name, category.Version)

	var model CategoryModel
	if err := model.Scan(row); err != nil {
		// Уникальность имени проверяет БД, в том числе при параллельных запросах.
		if errors.Is(err, core_postgres_pool.ErrViolatesUnique) {
			return domain.Category{}, fmt.Errorf("category name already exists: %w", core_errors.ErrConflict)
		}

		return domain.Category{}, fmt.Errorf("scan saved category: %w", err)
	}

	return modelToDomain(model), nil
}
