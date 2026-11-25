package resources_service

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"github.com/Yodau1/graphics-platform/internal/core/domain"
	core_errors "github.com/Yodau1/graphics-platform/internal/core/errors"
)

// GetPublishedResources возвращает страницу публичного каталога.
func (s *ResourcesService) GetPublishedResources(
	ctx context.Context,
	search string,
	categoryID *uuid.UUID,
	limit int,
	offset int,
) ([]domain.Resource, error) {
	if limit < 1 || limit > 100 {
		return nil, fmt.Errorf("limit must be between 1 and 100: %w", core_errors.ErrInvalidArgument)
	}

	if offset < 0 {
		return nil, fmt.Errorf("offset must be non-negative: %w", core_errors.ErrInvalidArgument)
	}

	search = strings.TrimSpace(search)

	resources, err := s.resourcesRepository.GetPublishedResources(ctx, search, categoryID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get published resources from repository: %w", err)
	}

	return resources, nil
}
