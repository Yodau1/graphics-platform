package users_service

import (
	"context"

	"github.com/google/uuid"

	"github.com/Yodau1/graphics-platform/internal/core/domain"
)

// UsersService — сервис пользователей с бизнес-логикой CRUD-операций.
type UsersService struct {
	usersRepository UsersRepository
}

// UsersRepository — интерфейс репозитория пользователей.
type UsersRepository interface {
	GetUser(
		ctx context.Context,
		id uuid.UUID,
	) (domain.User, error)
}

// NewUsersService создаёт сервис пользователей с внедрённым репозиторием
func NewUsersService(usersRepository UsersRepository) *UsersService {
	return &UsersService{
		usersRepository: usersRepository,
	}
}
