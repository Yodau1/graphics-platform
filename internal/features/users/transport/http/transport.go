package users_transport_http

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/Yodau1/graphics-platform/internal/core/domain"
	core_http_middleware "github.com/Yodau1/graphics-platform/internal/core/transport/http/middleware"
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
	PatchUser(
		ctx context.Context,
		id uuid.UUID,
		patch domain.UserPatch,
	) (domain.User, error)
	DeleteUser(
		ctx context.Context,
		id uuid.UUID,
	) error
	Login(
		ctx context.Context,
		email string,
		password string,
	) (string, time.Time, error)
	Authenticate(
		ctx context.Context,
		token string,
	) (domain.User, error)
	Logout(
		ctx context.Context,
		token string,
	) error
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
		{
			Method:  http.MethodPost,
			Path:    "/auth/login",
			Handler: h.Login,
		},
		{
			Method:     http.MethodGet,
			Path:       "/users/{id}",
			Handler:    h.GetUser,
			Middleware: []core_http_middleware.Middleware{h.authorizeUser(false)},
		},
		{
			Method:     http.MethodPatch,
			Path:       "/users/{id}",
			Handler:    h.PatchUser,
			Middleware: []core_http_middleware.Middleware{h.authorizeUser(false)},
		},
		{
			Method:     http.MethodDelete,
			Path:       "/users/{id}",
			Handler:    h.DeleteUser,
			Middleware: []core_http_middleware.Middleware{h.authorizeUser(true)},
		},
		{
			Method:  http.MethodPost,
			Path:    "/auth/logout",
			Handler: h.Logout,
		},
	}
}
