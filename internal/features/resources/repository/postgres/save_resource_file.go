package resources_postgres_repository

import (
	"context"
	"fmt"

	"github.com/Yodau1/graphics-platform/internal/core/domain"
	core_errors "github.com/Yodau1/graphics-platform/internal/core/errors"
)

// Добавляет файл и обновляет карточку одним sql запросом
func (r *ResourcesRepository) SaveResourceFile(
	ctx context.Context,
	resource domain.Resource,
	file domain.ResourceFile,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		WITH updated_resource AS (
			UPDATE graphics.resources
			SET status = 'draft',
			    version = version + 1,
			    updated_at = $4
			WHERE id = $1
			  AND author_id = $2
			  AND version = $3
			  AND status <> 'archived'
			RETURNING id
		)
		INSERT INTO graphics.resource_files (
			id, resource_id, kind, storage_key, original_name,
			format, size_bytes, created_at
		)
		SELECT $5, id, $6, $7, $8, $9, $10, $4
		FROM updated_resource;
	`

	tag, err := r.pool.Exec(
		ctx,
		query,
		resource.ID,
		resource.AuthorID,
		resource.Version,
		file.CreatedAt,
		file.ID,
		string(file.Kind),
		file.StorageKey,
		file.OriginalName,
		file.Format,
		file.SizeBytes,
	)
	if err != nil {
		return fmt.Errorf("save resource file metadata: %w", err)
	}

	// Ноль вставок означает, что карточка изменилась после чтения
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("resource changed concurrently: %w", core_errors.ErrConflict)
	}

	return nil
}
