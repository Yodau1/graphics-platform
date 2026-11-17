package domain

import (
	"fmt"
	"net/mail"
	"unicode/utf8"

	core_errors "github.com/Yodau1/graphics-platform/internal/core/errors"
	"github.com/google/uuid"
)

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

func CreateUser(name string, email string) User {
	return NewUser(
		uuid.New(),
		1,
		name,
		email,
		UserRoleUser,
		true,
	)
}

func (u *User) Validate() error {
	nameLength := utf8.RuneCountInString(u.Name)
	if nameLength < 3 || nameLength > 100 {
		return fmt.Errorf(
			"name must contain 3 to 100 characters: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	address, err := mail.ParseAddress(u.Email)
	if err != nil || address.Address != u.Email || utf8.RuneCountInString(u.Email) > 254 {
		return fmt.Errorf("invalid email address: %w", core_errors.ErrInvalidArgument)
	}
	return nil
}
