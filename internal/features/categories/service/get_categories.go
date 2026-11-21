package categories_service

import (
	"context"
	"fmt"

	"github.com/Yodau1/graphics-platform/internal/core/domain"
)

func (s *CategoriesService) GetCategories(ctx context.Context) ([]domain.Category, error) {
	categories, err := s.categoriesRepository.GetCategories(ctx)
	if err != nil {
		return nil, fmt.Errorf("get categories from repository: %w", err)
	}

	return categories, nil
}
