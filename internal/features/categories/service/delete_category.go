package categories_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s *CategoriesService) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	if err := s.categoriesRepository.DeleteCategory(ctx, id); err != nil {
		return fmt.Errorf("delete category: %w", err)
	}

	return nil
}
