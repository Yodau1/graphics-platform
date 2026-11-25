package resources_service

import (
	"context"
	"fmt"

	"github.com/Yodau1/graphics-platform/internal/core/domain"
	core_errors "github.com/Yodau1/graphics-platform/internal/core/errors"
	"github.com/google/uuid"
)

// Отправляем карточку автора на проверку
func (s *ResourcesService) SubmitResource(ctx context.Context, user domain.User, id uuid.UUID) (domain.Resource, error) {
	if user.Role != domain.UserRoleAuthor {
		return domain.Resource{}, fmt.Errorf("author role required: %w", core_errors.ErrForbidden)
	}

	resource, err := s.GetResource(ctx, user, id)
	if err != nil {
		return domain.Resource{}, fmt.Errorf("get resource fot submission: %w", err)
	}

	//Передаем репозиторию, для изменения статуса на "отправленный"
	submittedResource, err := s.resourcesRepository.SubmitResource(ctx, resource)
	if err != nil {
		return domain.Resource{}, fmt.Errorf("submit resource for review: %w", err)
	}

	return submittedResource, nil
}
