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

func (r *ResourcesRepository) HideResource(
	ctx context.Context,
	resourceID uuid.UUID,
	moderatorID uuid.UUID,
	version int64,
	reason string,
) (domain.Resource, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		WITH hidden_resource AS (
			UPDATE graphics.resources
			SET status = 'hidden',
				version = version + 1,
				updated_at = CURRENT_TIMESTAMP
			WHERE id = $1 AND version = $3 AND status = 'published'
			RETURNING id, version, author_id, category_id, title, description, status, created_at, updated_at
		),
		saved_event AS (
		INSERT INTO graphics.resource_events (
			id, resource_id, actor_id, action, details
		)
		SELECT $5, id, $2, 'hidden', jsonb_build_object('reason', $4::text)
		FROM hidden_resource
		)
		SELECT *
		FROM hidden_resource
	`

	row := r.pool.QueryRow(ctx, query, resourceID, moderatorID, version, reason, uuid.New())

	var resource ResourceModel
	if err := resource.Scan(row); err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.Resource{}, fmt.Errorf("resource is not published at the requested version: %w", core_errors.ErrConflict)
		}
		return domain.Resource{}, fmt.Errorf("scan hidden resource: %w", err)
	}

	return modelToDomain(resource), nil
}
