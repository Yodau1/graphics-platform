package resources_transport_http

//задаём JSON карточки и переносим значения из доменной модели

import (
	"time"

	"github.com/google/uuid"

	"github.com/Yodau1/graphics-platform/internal/core/domain"
)

type ResourceDTOResponse struct {
	ID          uuid.UUID             `json:"id"`
	Version     int64                 `json:"version"`
	AuthorID    uuid.UUID             `json:"author_id"`
	CategoryID  uuid.UUID             `json:"category_id"`
	Title       string                `json:"title"`
	Description string                `json:"description"`
	Status      domain.ResourceStatus `json:"status"`
	CreatedAt   time.Time             `json:"created_at"`
	UpdatedAt   *time.Time            `json:"updated_at"`
}

func resourceDTOFromDomain(resource domain.Resource) ResourceDTOResponse {
	return ResourceDTOResponse{
		ID:          resource.ID,
		Version:     resource.Version,
		AuthorID:    resource.AuthorID,
		CategoryID:  resource.CategoryID,
		Title:       resource.Title,
		Description: resource.Description,
		Status:      resource.Status,
		CreatedAt:   resource.CreatedAt,
		UpdatedAt:   resource.UpdatedAt,
	}
}

func resourceDTOsFromDomains(resources []domain.Resource) []ResourceDTOResponse {
	result := make([]ResourceDTOResponse, len(resources))

	for i, resource := range resources {
		result[i] = resourceDTOFromDomain(resource)
	}

	return result
}
