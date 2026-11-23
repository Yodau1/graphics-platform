package resources_service

import (
	"context"

	"github.com/google/uuid"

	"github.com/Yodau1/graphics-platform/internal/core/domain"
)

type ResourcesRepository interface {
	SaveResource(
		ctx context.Context,
		resource domain.Resource,
	) (domain.Resource, error)
	GetResource(
		ctx context.Context,
		id uuid.UUID,
	) (domain.Resource, error)
	GetResources(
		ctx context.Context,
		authorID uuid.UUID,
	) ([]domain.Resource, error)
	UpdateResource(
		ctx context.Context,
		resource domain.Resource,
	) (domain.Resource, error)
	SaveResourceFile(
		ctx context.Context,
		resource domain.Resource,
		file domain.ResourceFile,
	) error
	GetResourceFiles(
		ctx context.Context,
		resourceID uuid.UUID,
	) ([]domain.ResourceFile, error)
	DeleteResourceFile(
		ctx context.Context,
		resource domain.Resource,
		fileID uuid.UUID,
	) error
}

type ResourceFileStorage interface {
	SaveFile(ctx context.Context, key string, data []byte) error
	DeleteFile(key string) error
}

type ResourcesService struct {
	resourcesRepository ResourcesRepository
	fileStorage         ResourceFileStorage
}

func NewResourcesService(
	resourcesRepository ResourcesRepository,
	fileStorage ResourceFileStorage,
) *ResourcesService {
	return &ResourcesService{
		resourcesRepository: resourcesRepository,
		fileStorage:         fileStorage,
	}
}
