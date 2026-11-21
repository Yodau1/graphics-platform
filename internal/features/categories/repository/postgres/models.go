package categories_postgres_repository

import (
	"github.com/google/uuid"

	"github.com/Yodau1/graphics-platform/internal/core/domain"
	core_postgres_pool "github.com/Yodau1/graphics-platform/internal/core/repository/postgres/pool"
)

type CategoryModel struct {
	ID      uuid.UUID
	Name    string
	Version int64
}

func (m *CategoryModel) Scan(row core_postgres_pool.Row) error {
	return row.Scan(&m.ID, &m.Name, &m.Version)
}

func modelToDomain(model CategoryModel) domain.Category {
	return domain.NewCategory(model.ID, model.Version, model.Name)
}

func modelsToDomains(models []CategoryModel) []domain.Category {
	categories := make([]domain.Category, len(models))

	for i, model := range models {
		categories[i] = modelToDomain(model)
	}

	return categories
}
