package resources_service

import (
	"fmt"
	"strings"

	"github.com/Yodau1/graphics-platform/internal/core/domain"
	core_errors "github.com/Yodau1/graphics-platform/internal/core/errors"
)

// validateResourceFile проверяет имя, размер и формат файла.
func validateResourceFile(
	kind domain.ResourceFileKind,
	originalName string,
	data []byte,
) (string, string, error) {
	// TrimSpace убирает пробельные символы по краям имени.
	originalName = strings.TrimSpace(originalName)

	// len для строки и слайса байтов возвращает размер в байтах.
	if originalName == "" || len(originalName) > 255 ||
		len(data) == 0 || len(data) > MaxResourceFileSize {
		// Errorf с %w сохраняет ErrInvalidArgument для ответа HTTP 400.
		return "", "", fmt.Errorf("invalid file name or size: %w", core_errors.ErrInvalidArgument)
	}

	// Распознаём формат по содержимому файла.
	format := getResourceFileFormat(data)

	switch kind {
	case domain.ResourceFileKindOriginal:
		// Для оригинала подходит любой распознанный формат.
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
		return "", "", fmt.Errorf("unsupported format for file kind: %w", core_errors.ErrInvalidArgument)
	}

	return originalName, format, nil
}
