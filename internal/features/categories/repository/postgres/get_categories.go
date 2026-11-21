package categories_postgres_repository

import (
	"context"
	"fmt"

	"github.com/Yodau1/graphics-platform/internal/core/domain"
)

// Получение всех категорий, отсортированных по имени
func (r *CategoriesRepository) GetCategories(ctx context.Context) ([]domain.Category, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT id, name, version
		FROM graphics.categories
		ORDER BY name;
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query categories: %w", err)
	}
	defer rows.Close()

	var models []CategoryModel
	for rows.Next() {
		var model CategoryModel
		if err := model.Scan(rows); err != nil {
			return nil, fmt.Errorf("scan category: %w", err)
		}

		models = append(models, model)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate categories: %w", err)
	}

	return modelsToDomains(models), nil
}
