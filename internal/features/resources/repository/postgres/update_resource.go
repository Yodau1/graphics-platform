package resources_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Yodau1/graphics-platform/internal/core/domain"
	core_errors "github.com/Yodau1/graphics-platform/internal/core/errors"
	core_postgres_pool "github.com/Yodau1/graphics-platform/internal/core/repository/postgres/pool"
)

// UpdateResource сохраняет изменения, если версия карточки ещё актуальна.
func (r *ResourcesRepository) UpdateResource(
	ctx context.Context,
	resource domain.Resource,
) (domain.Resource, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		UPDATE graphics.resources
		SET category_id = $1,
		    title = $2,
		    description = $3,
		    status = $4,
		    updated_at = $5,
		    version = version + 1
		WHERE id = $6 AND version = $7
		RETURNING id, version, author_id, category_id, title, description,
		          status, created_at, updated_at;
	`

	row := r.pool.QueryRow(
		ctx,
		query,
		resource.CategoryID,
		resource.Title,
		resource.Description,
		string(resource.Status),
		resource.UpdatedAt,
		resource.ID,
		resource.Version,
	)

	var model ResourceModel
	if err := model.Scan(row); err != nil {
		switch {
		case errors.Is(err, core_postgres_pool.ErrNoRows):
			return domain.Resource{}, fmt.Errorf("resource changed concurrently: %w", core_errors.ErrConflict)
		case errors.Is(err, core_postgres_pool.ErrViolatesForeignKey):
			return domain.Resource{}, fmt.Errorf("resource category does not exist: %w", core_errors.ErrConflict)
		default:
			return domain.Resource{}, fmt.Errorf("scan updated resource: %w", err)
		}
	}

	return modelToDomain(model), nil
}
