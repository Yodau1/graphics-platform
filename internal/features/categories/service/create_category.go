package categories_service

import (
	"context"
	"fmt"

	"github.com/Yodau1/graphics-platform/internal/core/domain"
)

func (s *CategoriesService) CreateCategory(ctx context.Context, name string) (domain.Category, error) {
	category := domain.CreateCategory(name)
	if err := category.Validate(); err != nil {
		return domain.Category{}, fmt.Errorf("validate category name: %w", err)
	}

	saveCategory, err := s.categoriesRepository.SaveCategory(ctx, category)
	if err != nil {
		return domain.Category{}, fmt.Errorf("save category in repository: %w", err)
	}

	return saveCategory, nil
}
