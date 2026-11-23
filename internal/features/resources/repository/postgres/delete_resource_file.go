package resources_postgres_repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/Yodau1/graphics-platform/internal/core/domain"
	core_errors "github.com/Yodau1/graphics-platform/internal/core/errors"
)

// Снимает файл с карточки и увеличивает версию карточки
func (r *ResourcesRepository) DeleteResourceFile(
	ctx context.Context,
	resource domain.Resource,
	fileID uuid.UUID,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		WITH active_file AS (
			SELECT id
			FROM graphics.resource_files
			WHERE id = $4 AND resource_id = $1 AND retired_at IS NULL
		),
		updated_resource AS (
			UPDATE graphics.resources
			SET status = 'draft',
			    version = version + 1,
			    updated_at = $5
			WHERE id = $1
			  AND author_id = $2
			  AND version = $3
			  AND status <> 'archived'
			  AND EXISTS (SELECT 1 FROM active_file)
			RETURNING id
		),
		retired_file AS (
			UPDATE graphics.resource_files AS f
			SET retired_at = $5
			FROM updated_resource AS r, active_file AS a
			WHERE f.id = a.id AND f.resource_id = r.id
			RETURNING f.id
		)
		SELECT EXISTS (SELECT 1 FROM active_file),
		       EXISTS (SELECT 1 FROM retired_file);
	`

	updatedAt := time.Now().UTC()
	row := r.pool.QueryRow(
		ctx,
		query,
		resource.ID,
		resource.AuthorID,
		resource.Version,
		fileID,
		updatedAt,
	)

	var fileFound, fileRetired bool
	if err := row.Scan(&fileFound, &fileRetired); err != nil {
		return fmt.Errorf("delete resource file query: %w", err)
	}

	// Чужой, отсутствующий или уже снятый файл не входит в active_file
	if !fileFound {
		return fmt.Errorf("active resource file not found: %w", core_errors.ErrNotFound)
	}

	// Файл был найден, но карточка успела измениться после чтения
	if !fileRetired {
		return fmt.Errorf("resource changed concurrently: %w", core_errors.ErrConflict)
	}

	return nil
}
