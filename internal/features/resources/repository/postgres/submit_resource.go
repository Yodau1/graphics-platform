package resources_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Yodau1/graphics-platform/internal/core/domain"
	core_errors "github.com/Yodau1/graphics-platform/internal/core/errors"
	core_postgres_pool "github.com/Yodau1/graphics-platform/internal/core/repository/postgres/pool"
)

// Отправка черновика карточки на модерацию
func (r *ResourcesRepository) SubmitResource(ctx context.Context, resource domain.Resource) (domain.Resource, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		UPDATE graphics.resources AS r
		SET status = 'pending',
    		version = r.version + 1,
    		updated_at = CURRENT_TIMESTAMP
		WHERE r.id = $1
  		AND r.author_id = $2
  		AND r.version = $3
  		AND r.status = 'draft'
  		AND EXISTS (
      		SELECT 1
      		FROM graphics.resource_files AS f
      		WHERE f.resource_id = r.id
        		AND f.kind = 'original'
        		AND f.retired_at IS NULL
  		)
	RETURNING id, version, author_id, category_id, title, description,
          	status, created_at, updated_at;
	`

	row := r.pool.QueryRow(ctx, query, resource.ID, resource.AuthorID, resource.Version)

	var model ResourceModel

	if err := model.Scan(row); err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.Resource{}, fmt.Errorf("submission requires a current draft with an active original: %w", core_errors.ErrConflict)
		}
		return domain.Resource{}, fmt.Errorf("scan submitted resource: %w", err)
	}

	return modelToDomain(model), nil
}
