package resources_postgres_repository

import (
	"context"
	"fmt"

	"github.com/Yodau1/graphics-platform/internal/core/domain"
)

// Возвращает карточки, ожидающие проверки
func (r *ResourcesRepository) GetPendingResources(
	ctx context.Context,
) ([]domain.Resource, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT id, version, author_id, category_id, title, description,
		       status, created_at, updated_at
		FROM graphics.resources
		WHERE status = 'pending'
		ORDER BY updated_at ASC NULLS LAST, id ASC;
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query pending resources: %w", err)
	}
	defer rows.Close()

	resources := make([]domain.Resource, 0)

	for rows.Next() {
		var model ResourceModel

		if err := model.Scan(rows); err != nil {
			return nil, fmt.Errorf("scan pending resource: %w", err)
		}

		resources = append(resources, modelToDomain(model))
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate pending resources: %w", err)
	}

	return resources, nil
}
