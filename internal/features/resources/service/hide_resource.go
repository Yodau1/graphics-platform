package resources_service

import (
	"context"
	"fmt"
	"strings"

	"github.com/Yodau1/graphics-platform/internal/core/domain"
	core_errors "github.com/Yodau1/graphics-platform/internal/core/errors"
	"github.com/google/uuid"
)

func (s *ResourcesService) HideResource(
	ctx context.Context,
	moderator domain.User,
	resourceID uuid.UUID,
	version int64,
	reason string,
) (domain.Resource, error) {
	if moderator.Role != domain.UserRoleAdmin {
		return domain.Resource{}, fmt.Errorf("admin role required: %w", core_errors.ErrForbidden)
	}
	if version < 1 {
		return domain.Resource{}, fmt.Errorf("version must be greater than zero: %w", core_errors.ErrInvalidArgument)
	}

	reason = strings.TrimSpace(reason)
	if reason == "" {
		return domain.Resource{}, fmt.Errorf("hide reason is required: %w", core_errors.ErrInvalidArgument)
	}

	resource, err := s.resourcesRepository.HideResource(ctx, resourceID, moderator.ID, version, reason)
	if err != nil {
		return domain.Resource{}, fmt.Errorf("hide resource in repository: %w", err)
	}

	return resource, nil

}
