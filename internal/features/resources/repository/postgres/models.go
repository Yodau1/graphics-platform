package resources_postgres_repository

import (
	"time"

	"github.com/google/uuid"

	"github.com/Yodau1/graphics-platform/internal/core/domain"
	core_postgres_pool "github.com/Yodau1/graphics-platform/internal/core/repository/postgres/pool"
)

type ResourceModel struct {
	ID          uuid.UUID
	Version     int64
	AuthorID    uuid.UUID
	CategoryID  uuid.UUID
	Title       string
	Description string
	Status      string
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}

func (m *ResourceModel) Scan(row core_postgres_pool.Row) error {
	return row.Scan(
		&m.ID,
		&m.Version,
		&m.AuthorID,
		&m.CategoryID,
		&m.Title,
		&m.Description,
		&m.Status,
		&m.CreatedAt,
		&m.UpdatedAt,
	)
}

func modelToDomain(model ResourceModel) domain.Resource {
	return domain.NewResource(
		model.ID,
		model.Version,
		model.AuthorID,
		model.CategoryID,
		model.Title,
		model.Description,
		domain.ResourceStatus(model.Status),
		model.CreatedAt,
		nil,
	)
}
