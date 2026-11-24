package resources_postgres_repository

import (
	"context"
	"fmt"

	"github.com/Yodau1/graphics-platform/internal/core/domain"
	core_errors "github.com/Yodau1/graphics-platform/internal/core/errors"
)

// ReplaceResourceFile обновляет метаданные файла и версию карточки.
func (r *ResourcesRepository) ReplaceResourceFile(
	ctx context.Context,
	resource domain.Resource,
	oldFile domain.ResourceFile,
	newFile domain.ResourceFile,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		WITH updated_resource AS (
			UPDATE graphics.resources AS r
			SET status = 'draft',
			    version = r.version + 1,
			    updated_at = $6
			FROM graphics.resource_files AS f
			WHERE r.id = $1
			  AND r.author_id = $2
			  AND r.version = $3
			  AND r.status <> 'archived'
			  AND f.resource_id = r.id
			  AND f.id = $4
			  AND f.storage_key = $5
			  AND f.retired_at IS NULL
			RETURNING f.id AS file_id
		)
		UPDATE graphics.resource_files AS f
		SET storage_key = $7,
		    original_name = $8,
		    format = $9,
		    size_bytes = $10
		FROM updated_resource
		WHERE f.id = updated_resource.file_id;
	`

	tag, err := r.pool.Exec(
		ctx,
		query,
		resource.ID,
		resource.AuthorID,
		resource.Version,
		oldFile.ID,
		oldFile.StorageKey,
		newFile.CreatedAt,
		newFile.StorageKey,
		newFile.OriginalName,
		newFile.Format,
		newFile.SizeBytes,
	)
	if err != nil {
		return fmt.Errorf("replace resource file: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("resource or file changed concurrently: %w", core_errors.ErrConflict)
	}

	return nil
}
