package resources_transport_http

import (
	"net/http"

	"github.com/google/uuid"

	core_errors "github.com/Yodau1/graphics-platform/internal/core/errors"
	core_logger "github.com/Yodau1/graphics-platform/internal/core/logger"
	core_http_middleware "github.com/Yodau1/graphics-platform/internal/core/transport/http/middleware"
	core_http_request "github.com/Yodau1/graphics-platform/internal/core/transport/http/request"
	core_http_response "github.com/Yodau1/graphics-platform/internal/core/transport/http/response"
)

type CreateResourceRequest struct {
	CategoryID  uuid.UUID `json:"category_id"`
	Title       string    `json:"title" validate:"required"`
	Description string    `json:"description"`
}

type CreateResourceResponse ResourceDTOResponse

func (h *ResourcesHTTPHandler) CreateResource(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	user, ok := core_http_middleware.AuthenticatedUserFromContext(ctx)
	if !ok {
		responseHandler.ErrorResponse(core_errors.ErrUnauthorized, "authenticated user required")
		return
	}

	// Ограничиваем размер тела перед чтением JSON.
	r.Body = http.MaxBytesReader(rw, r.Body, 16*1024)

	var request CreateResourceRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")
		return
	}

	resource, err := h.resourcesService.CreateResource(
		ctx,
		user,
		request.CategoryID,
		request.Title,
		request.Description,
	)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create resource")
		return
	}

	response := CreateResourceResponse(resourceDTOFromDomain(resource))
	responseHandler.JSONResponse(response, http.StatusCreated)
}
