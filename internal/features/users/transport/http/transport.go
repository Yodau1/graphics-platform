package users_transport_http

import (
	"context"
	"net/http"

	"github.com/google/uuid"

	"github.com/Yodau1/graphics-platform/internal/core/domain"
	core_http_server "github.com/Yodau1/graphics-platform/internal/core/transport/http/server"
)

type UsersHTTPHandler struct {
	usersService UsersService
}

type UsersService interface {
	GetUser(
		ctx context.Context,
		id uuid.UUID,
	) (domain.User, error)
	CreateUser(
		ctx context.Context,
		name string,
		email string,
		password string,
	) (domain.User, error)
}

func NewUsersHTTPHandler(usersService UsersService) *UsersHTTPHandler {
	return &UsersHTTPHandler{
		usersService: usersService,
	}
}

func (h *UsersHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/users",
			Handler: h.CreateUser,
		},
	}
}
