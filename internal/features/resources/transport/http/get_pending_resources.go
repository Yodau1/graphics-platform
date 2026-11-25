package resources_transport_http

import (
	"net/http"

	core_logger "github.com/Yodau1/graphics-platform/internal/core/logger"
	core_http_response "github.com/Yodau1/graphics-platform/internal/core/transport/http/response"
)

type GetPendingResourcesResponse []ResourceDTOResponse

// Отправляет администратору очередь карточек на проверку.
func (h *ResourcesHTTPHandler) GetPendingResources(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	resources, err := h.resourcesService.GetPendingResources(ctx)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get pending resources")
		return
	}

	response := GetPendingResourcesResponse(resourceDTOsFromDomains(resources))
	responseHandler.JSONResponse(response, http.StatusOK)
}
