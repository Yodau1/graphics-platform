package users_service

import (
	"context"
	"time"

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
	SaveUser(
		ctx context.Context,
		user domain.User,
		passwordHash string,
	) (domain.User, error)
	UpdateUser(
		ctx context.Context,
		user domain.User,
	) (domain.User, error)
	DeleteUser(
		ctx context.Context,
		id uuid.UUID,
	) error
	GetUserCredentials(
		ctx context.Context,
		email string,
	) (uuid.UUID, string, bool, error)
	SaveUserSession(
		ctx context.Context,
		userID uuid.UUID,
		tokenHash []byte,
		expiresAt time.Time,
	) error
	GetUserBySession(ctx context.Context,
		tokenHash []byte,
	) (domain.User, error)
	DeleteUserSession(
		ctx context.Context,
		tokenHash []byte,
	) error
	BecomeAuthor(
		ctx context.Context,
		id uuid.UUID,
	) (domain.User, error)
	GetUsers(
		ctx context.Context,
		limit int,
		offset int,
	) ([]domain.User, error)
	UpdateUserStatus(
		ctx context.Context,
		id uuid.UUID,
		version int64,
		isActive bool,
	) (domain.User, error)
}

// NewUsersService создаёт сервис пользователей с внедрённым репозиторием
func NewUsersService(usersRepository UsersRepository) *UsersService {
	return &UsersService{
		usersRepository: usersRepository,
	}
}
