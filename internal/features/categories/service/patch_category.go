package categories_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/Yodau1/graphics-platform/internal/core/domain"
)

func (s *CategoriesService) PatchCategory(
	ctx context.Context,
	id uuid.UUID,
	patch domain.CategoryPatch,
) (domain.Category, error) {
	category, err := s.categoriesRepository.GetCategory(ctx, id)
	if err != nil {
		return domain.Category{}, fmt.Errorf("get category from repository: %w", err)
	}

	if err := category.ApplyPatch(patch); err != nil {
		return domain.Category{}, fmt.Errorf("apply category patch: %w", err)
	}

	patchedCategory, err := s.categoriesRepository.UpdateCategory(ctx, category)
	if err != nil {
		return domain.Category{}, fmt.Errorf("update category in repository: %w", err)
	}

	return patchedCategory, nil
}
