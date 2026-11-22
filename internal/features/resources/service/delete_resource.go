package resources_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/Yodau1/graphics-platform/internal/core/domain"
	core_errors "github.com/Yodau1/graphics-platform/internal/core/errors"
)

func (s *ResourcesService) DeleteResource(
	ctx context.Context,
	user domain.User,
	id uuid.UUID,
) error {
	if user.Role != domain.UserRoleAuthor {
		return fmt.Errorf("author role required: %w", core_errors.ErrForbidden)
	}

	resource, err := s.resourcesRepository.GetResource(ctx, id)
	if err != nil {
		return fmt.Errorf("get resource from repository: %w", err)
	}
	if resource.AuthorID != user.ID {
		return core_errors.ErrNotFound
	}

	if !resource.Archive() {
		return nil
	}

	if _, err := s.resourcesRepository.UpdateResource(ctx, resource); err != nil {
		return fmt.Errorf("archive resource in repository: %w", err)
	}

	return nil
}
