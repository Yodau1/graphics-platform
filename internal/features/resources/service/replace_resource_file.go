package resources_service

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"

	"github.com/Yodau1/graphics-platform/internal/core/domain"
	core_errors "github.com/Yodau1/graphics-platform/internal/core/errors"
)

func (s *ResourcesService) ReplaceResourceFile(
	ctx context.Context,
	user domain.User,
	resourceID uuid.UUID,
	fileID uuid.UUID,
	originalName string,
	data []byte,
) (domain.ResourceFile, error) {
	if user.Role != domain.UserRoleAuthor {
		return domain.ResourceFile{}, fmt.Errorf("author role required: %w", core_errors.ErrForbidden)
	}

	// Читает владельца, статус и текущую версию карточки из БД
	resource, err := s.resourcesRepository.GetResource(ctx, resourceID)
	if err != nil {
		return domain.ResourceFile{}, fmt.Errorf("get resource from repository: %w", err)
	}
	if resource.AuthorID != user.ID {
		return domain.ResourceFile{}, core_errors.ErrNotFound
	}
	if resource.Status == domain.ResourceStatusArchived {
		return domain.ResourceFile{}, fmt.Errorf("resource is archived: %w", core_errors.ErrConflict)
	}

	oldFile, err := s.resourcesRepository.GetResourceFile(ctx, resourceID, fileID)
	if err != nil {
		return domain.ResourceFile{}, fmt.Errorf("get resource file from repository: %w", err)
	}

	// Проверяем новое содержимое по назначению уже прикреплённого файла
	originalName, format, err := validateResourceFile(oldFile.Kind, originalName, data)
	if err != nil {
		return domain.ResourceFile{}, err
	}

	// Генерирует свободный ключ хранения для новых байтов.
	newFile := domain.CreateResourceFile(resourceID, oldFile.Kind, originalName, format, int64(len(data)))

	// Пишет содержимое под новым ключом. Старый файл в этот момент доступен.
	if err := s.fileStorage.SaveFile(ctx, newFile.StorageKey, data); err != nil {
		return domain.ResourceFile{}, fmt.Errorf("save replacement file: %w", err)
	}

	// Переключает ключ в БД и увеличивает версию карточки.
	if err := s.resourcesRepository.ReplaceResourceFile(ctx, resource, oldFile, newFile); err != nil {
		if errors.Is(err, core_errors.ErrConflict) {
			if cleanupErr := s.fileStorage.DeleteFile(newFile.StorageKey); cleanupErr != nil {
				log.Printf("delete unused replacement file for resource file %s: %v", fileID, cleanupErr)
			}
		}

		return domain.ResourceFile{}, fmt.Errorf("replace resource file in repository: %w", err)
	}

	// БД уже ссылается на новый ключ. Удаляем старое содержимое с диска.
	if err := s.fileStorage.DeleteFile(oldFile.StorageKey); err != nil {
		log.Printf("delete old content for resource file %s: %v", fileID, err)
	}

	newFile.ID = oldFile.ID
	newFile.CreatedAt = oldFile.CreatedAt

	return newFile, nil
}
