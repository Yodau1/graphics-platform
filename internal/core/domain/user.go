package domain

import (
	"fmt"
	"net/mail"
	"strings"
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

// UserPatch содержит изменения для частичного обновления пользователя (PATCH).
// Каждое поле обёрнуто в Nullable, чтобы различать «не передано» и «передано null».
type UserPatch struct {
	Name  Nullable[string]
	Email Nullable[string]
}

func NewUserPatch(
	name Nullable[string],
	email Nullable[string],
) UserPatch {
	return UserPatch{
		Name:  name,
		Email: email,
	}
}

func (p *UserPatch) Validate() error {
	if !p.Name.Set && !p.Email.Set {
		return fmt.Errorf("user patch must include name or email: %w", core_errors.ErrInvalidArgument)
	}

	if p.Name.Set && p.Name.Value == nil {
		return fmt.Errorf("name cannot be set to null: %w", core_errors.ErrInvalidArgument)
	}

	if p.Email.Set && p.Email.Value == nil {
		return fmt.Errorf("email cannot be set to null: %w", core_errors.ErrInvalidArgument)
	}

	return nil
}

func (u *User) ApplyPatch(patch UserPatch) error {
	if err := patch.Validate(); err != nil {
		return fmt.Errorf("validate user patch: %w", err)
	}

	tmp := *u

	if patch.Name.Set {
		tmp.Name = strings.TrimSpace(*patch.Name.Value)
	}
	if patch.Email.Set {
		tmp.Email = strings.ToLower(strings.TrimSpace(*patch.Email.Value))
	}
	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("validate patched user: %w", err)
	}

	*u = tmp

	return nil
}
