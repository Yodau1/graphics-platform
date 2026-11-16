package domain

import "github.com/google/uuid"

type UserRole string

const (
	UserRoleUser   UserRole = "user"
	UserRoleAuthor UserRole = "author"
	UserRoleAdmin  UserRole = "admin"
)

// Доменная сущность пользователя
type User struct {
	ID       uuid.UUID
	Name     string
	Email    string
	Role     UserRole
	IsActive bool
	Version  int64
}

// NewUser восстанавливает пользователя из сохранённых данных.
func NewUser(
	id uuid.UUID,
	version int64,
	name string,
	email string,
	role UserRole,
	isActive bool,
) User {
	return User{
		ID:       id,
		Version:  version,
		Name:     name,
		Email:    email,
		Role:     role,
		IsActive: isActive,
	}
}
