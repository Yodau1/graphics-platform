package resources_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Yodau1/graphics-platform/internal/core/domain"
	core_errors "github.com/Yodau1/graphics-platform/internal/core/errors"
	core_postgres_pool "github.com/Yodau1/graphics-platform/internal/core/repository/postgres/pool"
	"github.com/google/uuid"
)

func (r *ResourcesRepository) GetPublishedResource(ctx context.Context, id uuid.UUID) (domain.Resource, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT id, version, author_id, category_id, title, description,
		       status, created_at, updated_at
		FROM graphics.resources
		WHERE id = $1 AND status = 'published';
	`

	row := r.pool.QueryRow(ctx, query, id)

	var resource ResourceModel

	if err := resource.Scan(row); err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.Resource{}, fmt.Errorf("published resource not found: %w", core_errors.ErrNotFound)
		}
		return domain.Resource{}, fmt.Errorf("scan published resource: %w", err)
	}

	return modelToDomain(resource), nil
}
