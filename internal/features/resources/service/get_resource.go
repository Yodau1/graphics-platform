package resources_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/Yodau1/graphics-platform/internal/core/domain"
	core_errors "github.com/Yodau1/graphics-platform/internal/core/errors"
)

// GetResource возвращает карточку владельцу или администратору
func (s *ResourcesService) GetResource(
	ctx context.Context,
	user domain.User,
	id uuid.UUID,
) (domain.Resource, error) {
	resource, err := s.resourcesRepository.GetResource(ctx, id)
	if err != nil {
		return domain.Resource{}, fmt.Errorf("get resource from repository: %w", err)
	}

	if resource.AuthorID != user.ID && user.Role != domain.UserRoleAdmin {
		return domain.Resource{}, core_errors.ErrNotFound
	}

	return resource, nil
}
