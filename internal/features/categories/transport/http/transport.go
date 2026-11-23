package categories_transport_http

import (
	"context"
	"net/http"

	"github.com/google/uuid"

	"github.com/Yodau1/graphics-platform/internal/core/domain"
	core_http_middleware "github.com/Yodau1/graphics-platform/internal/core/transport/http/middleware"
	core_http_server "github.com/Yodau1/graphics-platform/internal/core/transport/http/server"
)

type CategoriesHTTPHandler struct {
	categoriesService    CategoriesService
	requireAuthenticated core_http_middleware.Middleware
}

type CategoriesService interface {
	CreateCategory(
		ctx context.Context,
		name string,
	) (domain.Category, error)
	GetCategories(
		ctx context.Context,
	) ([]domain.Category, error)
	GetCategory(
		ctx context.Context,
		id uuid.UUID,
	) (domain.Category, error)
	DeleteCategory(
		ctx context.Context,
		id uuid.UUID,
	) error
	PatchCategory(
		ctx context.Context,
		id uuid.UUID,
		patch domain.CategoryPatch,
	) (domain.Category, error)
}

func NewCategoriesHTTPHandler(
	categoriesService CategoriesService,
	requireAuthenticated core_http_middleware.Middleware,
) *CategoriesHTTPHandler {
	return &CategoriesHTTPHandler{
		categoriesService:    categoriesService,
		requireAuthenticated: requireAuthenticated,
	}
}

func (h *CategoriesHTTPHandler) Routes() []core_http_server.Route {
	adminMiddleware := []core_http_middleware.Middleware{
		h.requireAuthenticated,
		core_http_middleware.RequireAdmin(),
	}

	return []core_http_server.Route{
		{
			Method:     http.MethodPost,
			Path:       "/categories",
			Handler:    h.CreateCategory,
			Middleware: adminMiddleware,
		},
		{
			Method:  http.MethodGet,
			Path:    "/categories",
			Handler: h.GetCategories,
		},
		{
			Method:  http.MethodGet,
			Path:    "/categories/{id}",
			Handler: h.GetCategory,
		},
		{
			Method:     http.MethodPatch,
			Path:       "/categories/{id}",
			Handler:    h.PatchCategory,
			Middleware: adminMiddleware,
		},
		{
			Method:     http.MethodDelete,
			Path:       "/categories/{id}",
			Handler:    h.DeleteCategory,
			Middleware: adminMiddleware,
		},
	}
}
