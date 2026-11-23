package resources_transport_http

import (
	"time"

	"github.com/google/uuid"

	"github.com/Yodau1/graphics-platform/internal/core/domain"
)

// ResourceFileDTOResponse описывает метаданные, доступные клиенту.
type ResourceFileDTOResponse struct {
	ID           uuid.UUID               `json:"id"`
	ResourceID   uuid.UUID               `json:"resource_id"`
	Kind         domain.ResourceFileKind `json:"kind"`
	OriginalName string                  `json:"original_name"`
	Format       string                  `json:"format"`
	SizeBytes    int64                   `json:"size_bytes"`
	CreatedAt    time.Time               `json:"created_at"`
}

func resourceFileDTOFromDomain(file domain.ResourceFile) ResourceFileDTOResponse {
	return ResourceFileDTOResponse{
		ID:           file.ID,
		ResourceID:   file.ResourceID,
		Kind:         file.Kind,
		OriginalName: file.OriginalName,
		Format:       file.Format,
		SizeBytes:    file.SizeBytes,
		CreatedAt:    file.CreatedAt,
	}
}
