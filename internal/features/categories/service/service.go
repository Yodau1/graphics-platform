package categories_service

import (
	"context"

	"github.com/google/uuid"

	"github.com/Yodau1/graphics-platform/internal/core/domain"
)

type CategoriesService struct {
	categoriesRepository CategoriesRepository
}

// Требование сервиса к уровню репозитория
type CategoriesRepository interface {
	SaveCategory(
		ctx context.Context,
		category domain.Category,
	) (domain.Category, error)
	GetCategory(
		ctx context.Context,
		id uuid.UUID,
	) (domain.Category, error)
	GetCategories(
		ctx context.Context,
	) ([]domain.Category, error)
	DeleteCategory(
		ctx context.Context,
		id uuid.UUID,
	) error
	UpdateCategory(
		ctx context.Context,
		category domain.Category,
	) (domain.Category, error)
}

func NewCategoriesService(categoriesRepository CategoriesRepository) *CategoriesService {
	return &CategoriesService{
		categoriesRepository: categoriesRepository,
	}
}
