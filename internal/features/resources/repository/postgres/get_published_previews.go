package resources_postgres_repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/Yodau1/graphics-platform/internal/core/domain"
)

// GetPublishedPreviews возвращает метаданные действующих превью карточки
func (r *ResourcesRepository) GetPublishedPreviews(
	ctx context.Context,
	resourceID uuid.UUID,
) ([]domain.ResourceFile, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT f.id, f.resource_id, f.kind, f.original_name,
		       f.format, f.size_bytes, f.created_at
		FROM graphics.resource_files AS f
		JOIN graphics.resources AS r ON r.id = f.resource_id
		WHERE r.id = $1
		  AND r.status = 'published'
		  AND f.retired_at IS NULL
		  AND f.kind IN ('image', 'preview_3d')
		ORDER BY f.created_at ASC, f.id ASC;
	`

	rows, err := r.pool.Query(ctx, query, resourceID)
	if err != nil {
		return nil, fmt.Errorf("query published previews: %w", err)
	}
	defer rows.Close()

	previews := make([]domain.ResourceFile, 0)

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
			return nil, fmt.Errorf("scan published preview: %w", err)
		}

		file.Kind = domain.ResourceFileKind(kind)
		previews = append(previews, file)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate published previews: %w", err)
	}

	return previews, nil
}
