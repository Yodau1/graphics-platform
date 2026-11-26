package users_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/Yodau1/graphics-platform/internal/core/domain"
)

func (s *UsersService) PatchUserStatus(
	ctx context.Context,
	id uuid.UUID,
	isActive bool,
) (domain.User, error) {
	user, err := s.usersRepository.GetUser(ctx, id)
	if err != nil {
		return domain.User{}, fmt.Errorf("get user from repository: %w", err)
	}

	if user.IsActive == isActive {
		return user, nil
	}

	updatedUser, err := s.usersRepository.UpdateUserStatus(ctx, id, user.Version, isActive)
	if err != nil {
		return domain.User{}, fmt.Errorf("update user status in repository: %w", err)
	}

	return updatedUser, nil
}
