package resources_transport_http

import (
	"net/http"

	core_errors "github.com/Yodau1/graphics-platform/internal/core/errors"
	core_logger "github.com/Yodau1/graphics-platform/internal/core/logger"
	core_http_middleware "github.com/Yodau1/graphics-platform/internal/core/transport/http/middleware"
	core_http_response "github.com/Yodau1/graphics-platform/internal/core/transport/http/response"
)

type GetResourcesResponse []ResourceDTOResponse

func (h *ResourcesHTTPHandler) GetResources(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	user, ok := core_http_middleware.AuthenticatedUserFromContext(ctx)
	if !ok {
		responseHandler.ErrorResponse(core_errors.ErrUnauthorized, "authenticated user required")
		return
	}

	resources, err := h.resourcesService.GetResources(ctx, user)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get resources")
		return
	}

	response := GetResourcesResponse(resourceDTOsFromDomains(resources))
	responseHandler.JSONResponse(response, http.StatusOK)
}
