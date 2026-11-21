package categories_transport_http

import (
	"github.com/google/uuid"

	"github.com/Yodau1/graphics-platform/internal/core/domain"
)

type CategoryDTOResponse struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Version int64     `json:"version"`
}

func categoryDTOFromDomain(category domain.Category) CategoryDTOResponse {
	return CategoryDTOResponse{
		ID:      category.ID,
		Name:    category.Name,
		Version: category.Version,
	}
}

func categoryDTOsFromDomains(categories []domain.Category) []CategoryDTOResponse {
	result := make([]CategoryDTOResponse, len(categories))

	for i, category := range categories {
		result[i] = categoryDTOFromDomain(category)
	}

	return result
}
