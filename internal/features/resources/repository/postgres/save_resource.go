package resources_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Yodau1/graphics-platform/internal/core/domain"
	core_errors "github.com/Yodau1/graphics-platform/internal/core/errors"
	core_postgres_pool "github.com/Yodau1/graphics-platform/internal/core/repository/postgres/pool"
)

func (r *ResourcesRepository) SaveResource(
	ctx context.Context,
	resource domain.Resource,
) (domain.Resource, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		INSERT INTO graphics.resources (
			id, version, author_id, category_id, title, description,
			status, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, version, author_id, category_id, title, description,
		          status, created_at, updated_at;
	`

	row := r.pool.QueryRow(
		ctx,
		query,
		resource.ID,
		resource.Version,
		resource.AuthorID,
		resource.CategoryID,
		resource.Title,
		resource.Description,
		string(resource.Status),
		resource.CreatedAt,
		resource.UpdatedAt,
	)

	var model ResourceModel
	if err := model.Scan(row); err != nil {
		if errors.Is(err, core_postgres_pool.ErrViolatesForeignKey) {
			return domain.Resource{}, fmt.Errorf(
				"resource author or category does not exist: %w",
				core_errors.ErrConflict,
			)
		}

		return domain.Resource{}, fmt.Errorf("scan saved resource: %w", err)
	}

	resourceDomain := modelToDomain(model)

	return resourceDomain, nil
}
