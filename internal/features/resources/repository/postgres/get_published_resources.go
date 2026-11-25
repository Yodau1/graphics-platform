package resources_postgres_repository

import (
	"context"
	"fmt"

	"github.com/Yodau1/graphics-platform/internal/core/domain"
	"github.com/google/uuid"
)

// Поиск нескольких опубликованных ресурсов по названию, фильтрам и с пагинацией
func (r *ResourcesRepository) GetPublishedResources(
	ctx context.Context,
	search string,
	categoryID *uuid.UUID,
	limit int,
	offset int,
) ([]domain.Resource, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT id, version, author_id, category_id,
			title, description, status, created_at, updated_at
		FROM graphics.resources
		WHERE status = 'published'
			AND ($2::uuid IS NULL OR category_id = $2)
			AND ($1 = '' OR strpos(lower(title), lower($1)) > 0)
		ORDER BY updated_at DESC NULLS LAST, id DESC
		LIMIT $3 OFFSET $4;
	`

	rows, err := r.pool.Query(ctx, query, search, categoryID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("query published resources: %w", err)
	}
	//Освобождаем результат запроса
	defer rows.Close()

	resources := make([]domain.Resource, 0)

	for rows.Next() {
		var resource ResourceModel

		if err := resource.Scan(rows); err != nil {
			return nil, fmt.Errorf("scan published resource: %w", err)
		}

		resources = append(resources, modelToDomain(resource))
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate published resources: %w", err)
	}

	return resources, nil

}
