package resources_postgres_repository

import (
	"context"
	"fmt"

	"github.com/Yodau1/graphics-platform/internal/core/domain"
	"github.com/google/uuid"
)

func (r *ResourcesRepository) GetResources(ctx context.Context, authorID uuid.UUID) ([]domain.Resource, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT id, version, author_id, category_id, title, description,
		       status, created_at, updated_at
		FROM graphics.resources
		WHERE author_id = $1 AND status <> 'archived'
		ORDER BY created_at DESC, id DESC;
	`

	rows, err := r.pool.Query(ctx, query, authorID)
	if err != nil {
		return nil, fmt.Errorf("query resources: %w", err)
	}
	defer rows.Close()

	resources := make([]domain.Resource, 0)

	for rows.Next() {
		var resource ResourceModel
		if err := resource.Scan(rows); err != nil {
			return nil, fmt.Errorf("scan resources: %w", err)
		}

		resources = append(resources, modelToDomain(resource))
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate resources: %w", err)
	}

	return resources, nil

}
