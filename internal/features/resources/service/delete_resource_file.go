package resources_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/Yodau1/graphics-platform/internal/core/domain"
	core_errors "github.com/Yodau1/graphics-platform/internal/core/errors"
)

// DeleteResourceFile снимает файл с карточки её автора.
func (s *ResourcesService) DeleteResourceFile(
	ctx context.Context,
	user domain.User,
	resourceID uuid.UUID,
	fileID uuid.UUID,
) error {
	if user.Role != domain.UserRoleAuthor {
		return fmt.Errorf("author role required: %w", core_errors.ErrForbidden)
	}

	// Загружаем владельца, статус и текущую версию карточки.
	resource, err := s.resourcesRepository.GetResource(ctx, resourceID)
	if err != nil {
		return fmt.Errorf("get resource from repository: %w", err)
	}

	if resource.AuthorID != user.ID {
		return core_errors.ErrNotFound
	}

	if resource.Status == domain.ResourceStatusArchived {
		return fmt.Errorf("resource is archived: %w", core_errors.ErrConflict)
	}

	if err := s.resourcesRepository.DeleteResourceFile(ctx, resource, fileID); err != nil {
		return fmt.Errorf("delete resource file in repository: %w", err)
	}

	return nil
}
