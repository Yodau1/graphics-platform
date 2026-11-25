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
	UploadResourceFile(
		ctx context.Context,
		user domain.User,
		resourceID uuid.UUID,
		kind domain.ResourceFileKind,
		originalName string,
		data []byte,
	) (domain.ResourceFile, error)
	GetResourceFiles(
		ctx context.Context,
		user domain.User,
		resourceID uuid.UUID,
	) ([]domain.ResourceFile, error)

	DeleteResourceFile(
		ctx context.Context,
		user domain.User,
		resourceID uuid.UUID,
		fileID uuid.UUID,
	) error
	ReplaceResourceFile(
		ctx context.Context,
		user domain.User,
		resourceID uuid.UUID,
		fileID uuid.UUID,
		originalName string,
		data []byte,
	) (domain.ResourceFile, error)
	SubmitResource(
		ctx context.Context,
		user domain.User,
		resourceID uuid.UUID,
	) (domain.Resource, error)
	GetPendingResources(
		ctx context.Context,
	) ([]domain.Resource, error)
	ReviewResource(
		ctx context.Context,
		moderator domain.User,
		resourceID uuid.UUID,
		version int64,
		decision string,
		comment string,
	) (domain.Resource, error)
}

type SessionAuthenticator interface {
	Authenticate(ctx context.Context, token string) (domain.User, error)
}

type ResourcesHTTPHandler struct {
	resourcesService     ResourcesService
	requireAuthenticated core_http_middleware.Middleware
}

func NewResourcesHTTPHandler(
	resourcesService ResourcesService,
	requireAuthenticated core_http_middleware.Middleware,
) *ResourcesHTTPHandler {
	return &ResourcesHTTPHandler{
		resourcesService:     resourcesService,
		requireAuthenticated: requireAuthenticated,
	}
}

func (h *ResourcesHTTPHandler) Routes() []core_http_server.Route {
	// Оба маршрута работают с закрытыми карточками и требуют сессию
	auth := []core_http_middleware.Middleware{h.requireAuthenticated}

	adminMiddleware := []core_http_middleware.Middleware{
		//Определяем пользователя и кладем его в контекст запроса
		h.requireAuthenticated,
		//Проверяем роль админа у пользователя
		core_http_middleware.RequireAdmin(),
	}

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
		{
			Method:     http.MethodPost,
			Path:       "/resources/{id}/files",
			Handler:    h.UploadResourceFile,
			Middleware: auth,
		},
		{
			Method:     http.MethodGet,
			Path:       "/resources/{id}/files",
			Handler:    h.GetResourceFiles,
			Middleware: auth,
		},
		{
			Method:     http.MethodDelete,
			Path:       "/resources/{id}/files/{fileId}",
			Handler:    h.DeleteResourceFile,
			Middleware: auth,
		},
		{
			Method:     http.MethodPut,
			Path:       "/resources/{id}/files/{fileId}",
			Handler:    h.ReplaceResourceFile,
			Middleware: auth,
		},
		{
			Method:     http.MethodPost,
			Path:       "/resources/{id}/submit",
			Handler:    h.SubmitResource,
			Middleware: auth,
		},
		{
			Method:     http.MethodGet,
			Path:       "/moderation/resources",
			Handler:    h.GetPendingResources,
			Middleware: adminMiddleware,
		},
		{
			Method:     http.MethodPost,
			Path:       "/moderation/resources/{id}/review",
			Handler:    h.ReviewResource,
			Middleware: adminMiddleware,
		},
	}
}
