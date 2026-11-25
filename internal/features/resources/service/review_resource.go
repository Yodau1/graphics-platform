package resources_service

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"github.com/Yodau1/graphics-platform/internal/core/domain"
	core_errors "github.com/Yodau1/graphics-platform/internal/core/errors"
)

// Применяет решение администратора к проверенной версии карточки.
func (s *ResourcesService) ReviewResource(
	ctx context.Context,
	moderator domain.User,
	resourceID uuid.UUID,
	version int64,
	decision string,
	comment string,
) (domain.Resource, error) {
	if moderator.Role != domain.UserRoleAdmin {
		return domain.Resource{}, fmt.Errorf("administrator role required: %w", core_errors.ErrForbidden)
	}

	if version < 1 {
		return domain.Resource{}, fmt.Errorf("version must be greater than zero: %w", core_errors.ErrInvalidArgument)
	}

	// TrimSpace удаляет внешние пробелы. Комментарий из одних пробелов станет пустым.
	comment = strings.TrimSpace(comment)

	switch decision {
	case "approved":
	case "rejected":
		if comment == "" {
			return domain.Resource{}, fmt.Errorf("rejection requires a comment: %w", core_errors.ErrInvalidArgument)
		}
	default:
		return domain.Resource{}, fmt.Errorf("decision must be approved or rejected: %w", core_errors.ErrInvalidArgument)
	}

	reviewedResource, err := s.resourcesRepository.ReviewResource(
		ctx,
		resourceID,
		moderator.ID,
		version,
		decision,
		comment,
	)
	if err != nil {
		return domain.Resource{}, fmt.Errorf("review resource in repository: %w", err)
	}

	return reviewedResource, nil
}
