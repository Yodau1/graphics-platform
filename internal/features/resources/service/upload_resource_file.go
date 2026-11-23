package resources_service

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"image/jpeg"
	"image/png"
	"strings"

	"github.com/google/uuid"

	"github.com/Yodau1/graphics-platform/internal/core/domain"
	core_errors "github.com/Yodau1/graphics-platform/internal/core/errors"
)

const MaxResourceFileSize = 10 * 1024 * 1024

// Проверяет доступ, сохраняет файл и регистрирует его в БД
func (s *ResourcesService) UploadResourceFile(
	ctx context.Context,
	user domain.User,
	resourceID uuid.UUID,
	kind domain.ResourceFileKind,
	originalName string,
	data []byte,
) (domain.ResourceFile, error) {
	if user.Role != domain.UserRoleAuthor {
		return domain.ResourceFile{}, fmt.Errorf("author role required: %w", core_errors.ErrForbidden)
	}

	resource, err := s.resourcesRepository.GetResource(ctx, resourceID)
	if err != nil {
		return domain.ResourceFile{}, fmt.Errorf("get resource from repository: %w", err)
	}

	// Чужая карточка выглядит как отсутствующая
	if resource.AuthorID != user.ID {
		return domain.ResourceFile{}, core_errors.ErrNotFound
	}
	if resource.Status == domain.ResourceStatusArchived {
		return domain.ResourceFile{}, fmt.Errorf("resource is archived: %w", core_errors.ErrConflict)
	}

	originalName = strings.TrimSpace(originalName)
	if originalName == "" || len(originalName) > 255 ||
		len(data) == 0 || len(data) > MaxResourceFileSize {
		return domain.ResourceFile{}, fmt.Errorf("invalid file name or size: %w", core_errors.ErrInvalidArgument)
	}

	format := getResourceFileFormat(data)
	switch kind {
	// Для оригинала допустим любой из распознанных форматов
	case domain.ResourceFileKindOriginal:
	case domain.ResourceFileKindImage:
		if format != "png" && format != "jpeg" {
			format = ""
		}
	case domain.ResourceFileKindPreview3D:
		if format != "glb" {
			format = ""
		}
	default:
		format = ""
	}

	if format == "" {
		return domain.ResourceFile{}, fmt.Errorf("unsupported format for file kind: %w", core_errors.ErrInvalidArgument)
	}

	file := domain.CreateResourceFile(resource.ID, kind, originalName, format, int64(len(data)))

	// Сначала записываем содержимое, затем прикрепляем метаданные к карточке.
	if err := s.fileStorage.SaveFile(ctx, file.StorageKey, data); err != nil {
		return domain.ResourceFile{}, fmt.Errorf("save resource file: %w", err)
	}

	if err := s.resourcesRepository.SaveResourceFile(ctx, resource, file); err != nil {
		// При конфликте версий INSERT точно не выполнился
		if errors.Is(err, core_errors.ErrConflict) {
			if cleanupErr := s.fileStorage.DeleteFile(file.StorageKey); cleanupErr != nil {
				return domain.ResourceFile{}, fmt.Errorf("clean up rejected resource file: %w", cleanupErr)
			}
		}

		return domain.ResourceFile{}, fmt.Errorf("save resource file in repository: %w", err)
	}

	return file, nil
}

// Определяет формат по заголовку файла.
func getResourceFileFormat(data []byte) string {
	if _, err := png.DecodeConfig(bytes.NewReader(data)); err == nil {
		return "png"
	}

	if _, err := jpeg.DecodeConfig(bytes.NewReader(data)); err == nil {
		return "jpeg"
	}

	// Заголовок GLB 2.0 содержит сигнатуру, версию и полную длину файла.
	if len(data) >= 20 && string(data[:4]) == "glTF" &&
		binary.LittleEndian.Uint32(data[4:8]) == 2 &&
		binary.LittleEndian.Uint32(data[8:12]) == uint32(len(data)) {
		return "glb"
	}

	return ""
}
