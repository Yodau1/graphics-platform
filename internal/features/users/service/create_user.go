package users_service

import (
	"context"
	"fmt"
	"strings"

	"github.com/Yodau1/graphics-platform/internal/core/domain"
)

func (s *UsersService) CreateUser(
	ctx context.Context,
	name string,
	email string,
	password string,
) (domain.User, error) {
	user := domain.CreateUser(
		strings.TrimSpace(name),
		strings.ToLower(strings.TrimSpace(email)),
	)

	if err := user.Validate(); err != nil {
		return domain.User{}, fmt.Errorf("validate user domain: %w", err)
	}

	passwordHash, err := hashPassword(password)
	if err != nil {
		return domain.User{}, err
	}

	user, err = s.usersRepository.SaveUser(ctx, user, passwordHash)
	if err != nil {
		return domain.User{}, fmt.Errorf("save user in repository: %w", err)
	}
	return user, nil
}
