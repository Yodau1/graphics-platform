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

// Ввозвращает действующий файл предпросмотра опубликованной карточки.
func (r *ResourcesRepository) GetPublishedPreview(
	ctx context.Context,
	resourceID uuid.UUID,
	fileID uuid.UUID,
) (domain.ResourceFile, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT f.id, f.resource_id, f.kind, f.storage_key, f.original_name,
		       f.format, f.size_bytes, f.created_at
		FROM graphics.resource_files AS f
		JOIN graphics.resources AS r ON r.id = f.resource_id
		WHERE f.resource_id = $1
		  AND f.id = $2
		  AND r.status = 'published'
		  AND f.retired_at IS NULL
		  AND f.kind IN ('image', 'preview_3d');
	`

	row := r.pool.QueryRow(ctx, query, resourceID, fileID)

	var (
		file domain.ResourceFile
		kind string
	)

	if err := row.Scan(
		&file.ID,
		&file.ResourceID,
		&kind,
		&file.StorageKey,
		&file.OriginalName,
		&file.Format,
		&file.SizeBytes,
		&file.CreatedAt,
	); err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.ResourceFile{}, fmt.Errorf("published preview not found: %w", core_errors.ErrNotFound)
		}

		return domain.ResourceFile{}, fmt.Errorf("scan published preview: %w", err)
	}

	// Преобразуем строку из БД в доменный тип назначения файла
	file.Kind = domain.ResourceFileKind(kind)

	return file, nil
}
