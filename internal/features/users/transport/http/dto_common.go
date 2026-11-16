package users_transport_http

import (
	"github.com/google/uuid"

	"github.com/Yodau1/graphics-platform/internal/core/domain"
)

// UserDTOResponse — DTO для представления пользователя в API-ответе.
type UserDTOResponse struct {
	ID       uuid.UUID       `json:"id"`
	Name     string          `json:"name"`
	Email    string          `json:"email"`
	Role     domain.UserRole `json:"role"`
	IsActive bool            `json:"is_active"`
	Version  int64           `json:"version"`
}

// userDTOFromDomain конвертирует доменный объект User в DTO для HTTP-ответа.
func userDTOFromDomain(user domain.User) UserDTOResponse {
	return UserDTOResponse{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Role:     user.Role,
		IsActive: user.IsActive,
		Version:  user.Version,
	}
}
