package resources_transport_http

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"

	core_errors "github.com/Yodau1/graphics-platform/internal/core/errors"
	core_logger "github.com/Yodau1/graphics-platform/internal/core/logger"
	core_http_response "github.com/Yodau1/graphics-platform/internal/core/transport/http/response"
)

type GetPublishedResourcesResponse []ResourceDTOResponse

// Возвращает посетителю страницу каталога.
func (h *ResourcesHTTPHandler) GetPublishedResources(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	params := r.URL.Query()
	search := params.Get("search")
	limit, offset := 20, 0

	var categoryID *uuid.UUID
	if value := params.Get("category_id"); value != "" {
		id, err := uuid.Parse(value)
		if err != nil {
			responseHandler.ErrorResponse(core_errors.ErrInvalidArgument, "category_id must be a UUID")
			return
		}

		categoryID = &id
	}

	if value := params.Get("limit"); value != "" {
		parsedLimit, err := strconv.Atoi(value)
		if err != nil {
			responseHandler.ErrorResponse(core_errors.ErrInvalidArgument, "limit must be an integer")
			return
		}

		limit = parsedLimit
	}

	if value := params.Get("offset"); value != "" {
		parsedOffset, err := strconv.Atoi(value)
		if err != nil {
			responseHandler.ErrorResponse(core_errors.ErrInvalidArgument, "offset must be an integer")
			return
		}

		offset = parsedOffset
	}

	resources, err := h.resourcesService.GetPublishedResources(ctx, search, categoryID, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(fmt.Errorf("get catalog page: %w", err), "failed to get published resources")
		return
	}

	response := GetPublishedResourcesResponse(resourceDTOsFromDomains(resources))
	responseHandler.JSONResponse(response, http.StatusOK)
}
