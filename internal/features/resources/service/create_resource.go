package resources_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/Yodau1/graphics-platform/internal/core/domain"
	core_errors "github.com/Yodau1/graphics-platform/internal/core/errors"
)

// CreateResource создаёт черновик от имени автора
func (s *ResourcesService) CreateResource(
	ctx context.Context,
	user domain.User,
	categoryID uuid.UUID,
	title string,
	description string,
) (domain.Resource, error) {
	if user.Role != domain.UserRoleAuthor {
		return domain.Resource{}, fmt.Errorf("author role required: %w", core_errors.ErrForbidden)
	}

	resource := domain.CreateResource(user.ID, categoryID, title, description)
	if err := resource.Validate(); err != nil {
		return domain.Resource{}, fmt.Errorf("validate resource domain: %w", err)
	}

	savedResource, err := s.resourcesRepository.SaveResource(ctx, resource)
	if err != nil {
		return domain.Resource{}, fmt.Errorf("save resource in repository: %w", err)
	}

	return savedResource, nil
}
