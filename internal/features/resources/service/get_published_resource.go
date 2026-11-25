package resources_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/Yodau1/graphics-platform/internal/core/domain"
)

// Возвращает опубликованную карточку по ID и ее превью.
func (s *ResourcesService) GetPublishedResource(
	ctx context.Context,
	id uuid.UUID,
) (domain.Resource, []domain.ResourceFile, error) {
	resource, err := s.resourcesRepository.GetPublishedResource(ctx, id)
	if err != nil {
		return domain.Resource{}, nil, fmt.Errorf("get published resource from repository: %w", err)
	}

	// GetPublishedPreviews читает метаданные изображений и GLB-превью этой карточки.
	previews, err := s.resourcesRepository.GetPublishedPreviews(ctx, resource.ID)
	if err != nil {
		return domain.Resource{}, nil, fmt.Errorf("get published previews from repository: %w", err)
	}

	return resource, previews, nil
}
