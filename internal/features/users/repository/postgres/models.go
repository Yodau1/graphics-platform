package users_postgres_repository

import (
	"github.com/Yodau1/graphics-platform/internal/core/domain"
	core_postgres_pool "github.com/Yodau1/graphics-platform/internal/core/repository/postgres/pool"
	"github.com/google/uuid"
)

// Модель пользователя для получения из репозитория
type userModel struct {
	ID       uuid.UUID
	Name     string
	Email    string
	Role     string
	IsActive bool
	Version  int64
}

// Scan записывает поля результата row в userModel
func (m *userModel) Scan(row core_postgres_pool.Row) error {
	return row.Scan(
		&m.ID,
		&m.Name,
		&m.Email,
		&m.Role,
		&m.IsActive,
		&m.Version,
	)
}

func modelToDomain(model userModel) domain.User {
	return domain.NewUser(
		model.ID,
		model.Version,
		model.Name,
		model.Email,
		domain.UserRole(model.Role),
		model.IsActive,
	)
}
