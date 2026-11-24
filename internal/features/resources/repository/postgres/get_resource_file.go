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

// Возвращает действующий файл указанной карточки.
func (r *ResourcesRepository) GetResourceFile(ctx context.Context, resourceID uuid.UUID, fileID uuid.UUID) (domain.ResourceFile, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel() // Освобождаем таймер при выходе из метода.

	query := `
		SELECT id, resource_id, kind, storage_key, original_name,
		       format, size_bytes, created_at
		FROM graphics.resource_files
		WHERE resource_id = $1 AND id = $2 AND retired_at IS NULL;
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
			return domain.ResourceFile{}, fmt.Errorf("active resource file not found: %w", core_errors.ErrNotFound)
		}

		return domain.ResourceFile{}, fmt.Errorf("scan resource file: %w", err)
	}

	file.Kind = domain.ResourceFileKind(kind)

	return file, nil
}
