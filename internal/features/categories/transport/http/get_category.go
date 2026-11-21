package categories_transport_http

import (
	"net/http"

	core_logger "github.com/Yodau1/graphics-platform/internal/core/logger"
	core_http_request "github.com/Yodau1/graphics-platform/internal/core/transport/http/request"
	core_http_response "github.com/Yodau1/graphics-platform/internal/core/transport/http/response"
)

type GetCategoryResponse CategoryDTOResponse

func (h *CategoriesHTTPHandler) GetCategory(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	categoryID, err := core_http_request.GetUUIDPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get categoryID path value")
		return
	}

	category, err := h.categoriesService.GetCategory(ctx, categoryID)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get category")
		return
	}

	response := GetCategoryResponse(categoryDTOFromDomain(category))
	responseHandler.JSONResponse(response, http.StatusOK)
}
