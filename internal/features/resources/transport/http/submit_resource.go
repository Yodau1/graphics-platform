package resources_transport_http

import (
	"net/http"

	core_errors "github.com/Yodau1/graphics-platform/internal/core/errors"
	core_logger "github.com/Yodau1/graphics-platform/internal/core/logger"
	core_http_middleware "github.com/Yodau1/graphics-platform/internal/core/transport/http/middleware"
	core_http_request "github.com/Yodau1/graphics-platform/internal/core/transport/http/request"
	core_http_response "github.com/Yodau1/graphics-platform/internal/core/transport/http/response"
)

type SubmitResourceResponse ResourceDTOResponse

func (h *ResourcesHTTPHandler) SubmitResource(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	// AuthenticatedUserFromContext читает пользователя, которого уже определило middleware
	user, ok := core_http_middleware.AuthenticatedUserFromContext(ctx)
	if !ok {
		responseHandler.ErrorResponse(core_errors.ErrUnauthorized, "authenticated user required")
		return
	}

	resourceID, err := core_http_request.GetUUIDPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get resourceID path value")
		return
	}

	resource, err := h.resourcesService.SubmitResource(ctx, user, resourceID)
	if err != nil {
		// ErrorResponse выбирает статус по типу ошибки сервиса.
		responseHandler.ErrorResponse(err, "failed to submit resource for review")
		return
	}

	response := SubmitResourceResponse(resourceDTOFromDomain(resource))
	responseHandler.JSONResponse(response, http.StatusOK)
}
