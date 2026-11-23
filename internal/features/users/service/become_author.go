package users_service

import (
	"context"
	"fmt"

	"github.com/Yodau1/graphics-platform/internal/core/domain"
)

func (s *UsersService) BecomeAuthor(ctx context.Context, token string) (domain.User, error) {

	user, err := s.Authenticate(ctx, token)
	if err != nil {
		return domain.User{}, fmt.Errorf("authenticate user: %w", err)
	}

	updatedUser, err := s.usersRepository.BecomeAuthor(ctx, user.ID)
	if err != nil {
		return domain.User{}, fmt.Errorf("updated user role: %w", err)
	}

	return updatedUser, nil

}
