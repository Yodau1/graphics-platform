package resources_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/Yodau1/graphics-platform/internal/core/domain"
	core_errors "github.com/Yodau1/graphics-platform/internal/core/errors"
	core_postgres_pool "github.com/Yodau1/graphics-platform/internal/core/repository/postgres/pool"
)

func (r *ResourcesRepository) GetResource(
	ctx context.Context,
	id uuid.UUID,
) (domain.Resource, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT id, version, author_id, category_id, title, description,
		       status, created_at, updated_at
		FROM graphics.resources
		WHERE id = $1;
	`

	row := r.pool.QueryRow(ctx, query, id)

	var model ResourceModel
	if err := model.Scan(row); err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.Resource{}, fmt.Errorf("resource with id='%s': %w", id, core_errors.ErrNotFound)
		}

		return domain.Resource{}, fmt.Errorf("scan resource: %w", err)
	}

	resourceDomain := modelToDomain(model)

	return resourceDomain, nil
}
