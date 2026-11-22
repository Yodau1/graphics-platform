package resources_transport_http

import (
	"context"
	"net/http"

	"github.com/google/uuid"

	"github.com/Yodau1/graphics-platform/internal/core/domain"
	core_http_middleware "github.com/Yodau1/graphics-platform/internal/core/transport/http/middleware"
	core_http_server "github.com/Yodau1/graphics-platform/internal/core/transport/http/server"
)

type ResourcesService interface {
	CreateResource(
		ctx context.Context,
		user domain.User,
		categoryID uuid.UUID,
		title string,
		description string,
	) (domain.Resource, error)
	GetResource(
		ctx context.Context,
		user domain.User,
		id uuid.UUID,
	) (domain.Resource, error)
	GetResources(
		ctx context.Context,
		user domain.User,
	) ([]domain.Resource, error)
	PatchResource(
		ctx context.Context,
		user domain.User,
		id uuid.UUID,
		patch domain.ResourcePatch,
	) (domain.Resource, error)
	DeleteResource(
		ctx context.Context,
		user domain.User,
		id uuid.UUID,
	) error
}

type SessionAuthenticator interface {
	Authenticate(ctx context.Context, token string) (domain.User, error)
}

type ResourcesHTTPHandler struct {
	resourcesService ResourcesService
	authenticator    SessionAuthenticator
}

func NewResourcesHTTPHandler(
	resourcesService ResourcesService,
	authenticator SessionAuthenticator,
) *ResourcesHTTPHandler {
	return &ResourcesHTTPHandler{
		resourcesService: resourcesService,
		authenticator:    authenticator,
	}
}

func (h *ResourcesHTTPHandler) Routes() []core_http_server.Route {
	// Оба маршрута работают с закрытыми карточками и требуют сессию
	auth := []core_http_middleware.Middleware{h.requireAuthenticatedUser()}

	return []core_http_server.Route{
		{
			Method:     http.MethodPost,
			Path:       "/resources",
			Handler:    h.CreateResource,
			Middleware: auth,
		},
		{
			Method:     http.MethodGet,
			Path:       "/resources/{id}",
			Handler:    h.GetResource,
			Middleware: auth,
		},
		{
			Method:     http.MethodGet,
			Path:       "/resources",
			Handler:    h.GetResources,
			Middleware: auth,
		},
		{
			Method:     http.MethodPatch,
			Path:       "/resources/{id}",
			Handler:    h.PatchResource,
			Middleware: auth,
		},
		{
			Method:     http.MethodDelete,
			Path:       "/resources/{id}",
			Handler:    h.DeleteResource,
			Middleware: auth,
		},
	}
}
