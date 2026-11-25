package resources_transport_http

import (
	"net/http"

	core_errors "github.com/Yodau1/graphics-platform/internal/core/errors"
	core_logger "github.com/Yodau1/graphics-platform/internal/core/logger"
	core_http_middleware "github.com/Yodau1/graphics-platform/internal/core/transport/http/middleware"
	core_http_request "github.com/Yodau1/graphics-platform/internal/core/transport/http/request"
	core_http_response "github.com/Yodau1/graphics-platform/internal/core/transport/http/response"
)

// ReviewResourceRequest содержит решение по конкретной версии карточки.
type ReviewResourceRequest struct {
	Version  int64  `json:"version" validate:"required"`
	Decision string `json:"decision" validate:"required"`
	Comment  string `json:"comment"`
}

type ReviewResourceResponse ResourceDTOResponse

// Принимает решение администратора по карточке из очереди.
func (h *ResourcesHTTPHandler) ReviewResource(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	moderator, ok := core_http_middleware.AuthenticatedUserFromContext(ctx)
	if !ok {
		// ErrorResponse преобразует ErrUnauthorized в HTTP 401.
		responseHandler.ErrorResponse(core_errors.ErrUnauthorized, "authenticated user required")
		return
	}

	resourceID, err := core_http_request.GetUUIDPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get resourceID path value")
		return
	}

	r.Body = http.MaxBytesReader(rw, r.Body, 16*1024)

	var request ReviewResourceRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")
		return
	}

	resource, err := h.resourcesService.ReviewResource(
		ctx,
		moderator,
		resourceID,
		request.Version,
		request.Decision,
		request.Comment,
	)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to review resource")
		return
	}

	response := ReviewResourceResponse(resourceDTOFromDomain(resource))
	responseHandler.JSONResponse(response, http.StatusOK)
}
