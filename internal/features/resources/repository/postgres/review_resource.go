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

// Сохраняет решение модератора и возвращает обновлённую карточку.
func (r *ResourcesRepository) ReviewResource(
	ctx context.Context,
	resourceID uuid.UUID,
	moderatorID uuid.UUID,
	version int64,
	decision string,
	comment string,
) (domain.Resource, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		WITH reviewed_resource AS (
			UPDATE graphics.resources
			SET status = CASE $4::text
			        WHEN 'approved' THEN 'published'
			        WHEN 'rejected' THEN 'rejected'
			    END,
			    version = version + 1,
			    updated_at = CURRENT_TIMESTAMP
			WHERE id = $1 AND version = $3 AND status = 'pending'
			RETURNING id, version, author_id, category_id, title, description,
			          status, created_at, updated_at
		),
		saved_review AS (
			INSERT INTO graphics.moderation_reviews (
				id, resource_id, moderator_id, reviewed_version, decision, comment
			)
			SELECT $6, id, $2, $3, $4, $5
			FROM reviewed_resource
			RETURNING resource_id
		)
		SELECT r.*
		FROM reviewed_resource AS r
		JOIN saved_review AS s ON s.resource_id = r.id;
	`

	reviewID := uuid.New()

	row := r.pool.QueryRow(
		ctx,
		query,
		resourceID,
		moderatorID,
		version,
		decision,
		comment,
		reviewID,
	)

	var model ResourceModel
	if err := model.Scan(row); err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.Resource{}, fmt.Errorf("resource is no longer pending at the reviewed version: %w", core_errors.ErrConflict)
		}

		return domain.Resource{}, fmt.Errorf("scan reviewed resource: %w", err)
	}

	return modelToDomain(model), nil
}
