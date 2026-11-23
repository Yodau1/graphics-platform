package resources_postgres_repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/Yodau1/graphics-platform/internal/core/domain"
)

// Возвращает действующие файлы карточки
func (r *ResourcesRepository) GetResourceFiles(
	ctx context.Context,
	resourceID uuid.UUID,
) ([]domain.ResourceFile, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT id, resource_id, kind, original_name, format, size_bytes, created_at
		FROM graphics.resource_files
		WHERE resource_id = $1 AND retired_at IS NULL
		ORDER BY created_at, id;
	`

	rows, err := r.pool.Query(ctx, query, resourceID)
	if err != nil {
		return nil, fmt.Errorf("query resource files: %w", err)
	}
	defer rows.Close()

	files := make([]domain.ResourceFile, 0)
	for rows.Next() {
		var (
			file domain.ResourceFile
			kind string
		)

		if err := rows.Scan(
			&file.ID,
			&file.ResourceID,
			&kind,
			&file.OriginalName,
			&file.Format,
			&file.SizeBytes,
			&file.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan resource file: %w", err)
		}

		file.Kind = domain.ResourceFileKind(kind)
		files = append(files, file)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate resource files: %w", err)
	}

	return files, nil
}
