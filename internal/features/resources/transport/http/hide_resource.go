package resources_transport_http

import (
	"net/http"

	core_errors "github.com/Yodau1/graphics-platform/internal/core/errors"
	core_logger "github.com/Yodau1/graphics-platform/internal/core/logger"
	core_http_middleware "github.com/Yodau1/graphics-platform/internal/core/transport/http/middleware"
	core_http_request "github.com/Yodau1/graphics-platform/internal/core/transport/http/request"
	core_http_response "github.com/Yodau1/graphics-platform/internal/core/transport/http/response"
)

type HideResourceRequest struct {
	Version int64  `json:"version" validate:"required,min=1"`
	Reason  string `json:"reason" validate:"required"`
}

type HideResourceResponse ResourceDTOResponse

func (h *ResourcesHTTPHandler) HideResource(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	moderator, ok := core_http_middleware.AuthenticatedUserFromContext(ctx)
	if !ok {
		responseHandler.ErrorResponse(core_errors.ErrUnauthorized, "authenticated user required")
		return
	}

	resourceID, err := core_http_request.GetUUIDPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get resourceID path value")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 16*1024)

	var request HideResourceRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")
		return
	}

	resource, err := h.resourcesService.HideResource(
		ctx,
		moderator,
		resourceID,
		request.Version,
		request.Reason,
	)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to hide resource")
		return
	}

	response := HideResourceResponse(resourceDTOFromDomain(resource))
	responseHandler.JSONResponse(response, http.StatusOK)
}
