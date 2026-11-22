package resources_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/Yodau1/graphics-platform/internal/core/domain"
	core_errors "github.com/Yodau1/graphics-platform/internal/core/errors"
)

func (s *ResourcesService) PatchResource(
	ctx context.Context,
	user domain.User,
	id uuid.UUID,
	patch domain.ResourcePatch,
) (domain.Resource, error) {
	if user.Role != domain.UserRoleAuthor {
		return domain.Resource{}, fmt.Errorf("author role required: %w", core_errors.ErrForbidden)
	}

	resource, err := s.resourcesRepository.GetResource(ctx, id)
	if err != nil {
		return domain.Resource{}, fmt.Errorf("get resource from repository: %w", err)
	}
	if resource.AuthorID != user.ID {
		return domain.Resource{}, core_errors.ErrNotFound
	}

	if err := resource.ApplyPatch(patch); err != nil {
		return domain.Resource{}, fmt.Errorf("apply resource patch: %w", err)
	}

	patchedResource, err := s.resourcesRepository.UpdateResource(ctx, resource)
	if err != nil {
		return domain.Resource{}, fmt.Errorf("update resource in repository: %w", err)
	}

	return patchedResource, nil
}
