package categories_transport_http

import (
	"net/http"

	core_logger "github.com/Yodau1/graphics-platform/internal/core/logger"
	core_http_request "github.com/Yodau1/graphics-platform/internal/core/transport/http/request"
	core_http_response "github.com/Yodau1/graphics-platform/internal/core/transport/http/response"
)

type CreateCategoryRequest struct {
	Name string `json:"name" validate:"required"`
}

type CreateCategoryResponse CategoryDTOResponse

func (h *CategoriesHTTPHandler) CreateCategory(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	r.Body = http.MaxBytesReader(rw, r.Body, 16*1024)

	var request CreateCategoryRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")
		return
	}

	category, err := h.categoriesService.CreateCategory(ctx, request.Name)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create category")
		return
	}

	response := CreateCategoryResponse(categoryDTOFromDomain(category))
	responseHandler.JSONResponse(response, http.StatusCreated)
}
