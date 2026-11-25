package resources_service

import (
	"context"
	"fmt"

	"github.com/Yodau1/graphics-platform/internal/core/domain"
)

func (s *ResourcesService) GetPendingResources(ctx context.Context) ([]domain.Resource, error) {
	resources, err := s.resourcesRepository.GetPendingResources(ctx)
	if err != nil {
		return nil, fmt.Errorf("get pending resources from repository: %w", err)
	}

	return resources, nil
}
