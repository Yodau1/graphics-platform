package resources_service

import (
	"context"
	"fmt"

	"github.com/Yodau1/graphics-platform/internal/core/domain"
)

func (s *ResourcesService) GetResources(
	ctx context.Context,
	user domain.User,
) ([]domain.Resource, error) {
	resources, err := s.resourcesRepository.GetResources(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("get resources from repository: %w", err)
	}

	return resources, nil
}
