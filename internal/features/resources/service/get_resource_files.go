package resources_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/Yodau1/graphics-platform/internal/core/domain"
)

// Возвращает файлы карточки, которую пользователь вправе читать
func (s *ResourcesService) GetResourceFiles(
	ctx context.Context,
	user domain.User,
	resourceID uuid.UUID,
) ([]domain.ResourceFile, error) {
	// Проверяем доступ к карточке до чтения файлов.
	resource, err := s.GetResource(ctx, user, resourceID)
	if err != nil {
		return nil, fmt.Errorf("get resource: %w", err)
	}

	files, err := s.resourcesRepository.GetResourceFiles(ctx, resource.ID)
	if err != nil {
		return nil, fmt.Errorf("get resource files from repository: %w", err)
	}

	return files, nil
}
