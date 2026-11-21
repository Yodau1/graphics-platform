package categories_transport_http

import (
	"net/http"

	core_logger "github.com/Yodau1/graphics-platform/internal/core/logger"
	core_http_response "github.com/Yodau1/graphics-platform/internal/core/transport/http/response"
)

type GetCategoriesResponse []CategoryDTOResponse

func (h *CategoriesHTTPHandler) GetCategories(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	categories, err := h.categoriesService.GetCategories(ctx)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get categories")

		return
	}

	// Переводим доменные сущности в публичный формат ответа.
	response := GetCategoriesResponse(categoryDTOsFromDomains(categories))
	responseHandler.JSONResponse(response, http.StatusOK)
}
