package categories_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/Yodau1/graphics-platform/internal/core/domain"
)

func (s *CategoriesService) GetCategory(ctx context.Context, id uuid.UUID) (domain.Category, error) {
	category, err := s.categoriesRepository.GetCategory(ctx, id)
	if err != nil {
		return domain.Category{}, fmt.Errorf("get category from repository: %w", err)
	}

	return category, nil
}
