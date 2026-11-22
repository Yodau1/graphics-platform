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
}

type ResourcesService struct {
	resourcesRepository ResourcesRepository
}

func NewResourcesService(resourcesRepository ResourcesRepository) *ResourcesService {
	return &ResourcesService{
		resourcesRepository: resourcesRepository,
	}
}
