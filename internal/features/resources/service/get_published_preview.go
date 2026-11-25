package resources_service

import (
	"context"
	"fmt"
	"io"

	"github.com/google/uuid"

	"github.com/Yodau1/graphics-platform/internal/core/domain"
)

func (s *ResourcesService) GetPublishedPreview(
	ctx context.Context,
	resourceID uuid.UUID,
	fileID uuid.UUID,
) (domain.ResourceFile, io.ReadSeekCloser, error) {
	file, err := s.resourcesRepository.GetPublishedPreview(ctx, resourceID, fileID)
	if err != nil {
		return domain.ResourceFile{}, nil, fmt.Errorf("get published preview: %w", err)
	}

	content, err := s.fileStorage.OpenFile(ctx, file.StorageKey)
	if err != nil {
		return domain.ResourceFile{}, nil, fmt.Errorf("open published preview: %w", err)
	}

	return file, content, nil
}
