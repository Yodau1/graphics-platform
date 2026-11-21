package categories_transport_http

import (
	"net/http"

	"github.com/Yodau1/graphics-platform/internal/core/domain"
	core_logger "github.com/Yodau1/graphics-platform/internal/core/logger"
	core_http_request "github.com/Yodau1/graphics-platform/internal/core/transport/http/request"
	core_http_response "github.com/Yodau1/graphics-platform/internal/core/transport/http/response"
	core_http_types "github.com/Yodau1/graphics-platform/internal/core/transport/http/types"
)

type PatchCategoryRequest struct {
	Name core_http_types.Nullable[string] `json:"name" swaggertype:"string"`
}

// Validate проверяет структуру патча до обращения к сервису.
func (r *PatchCategoryRequest) Validate() error {
	patch := categoryPatchFromRequest(*r)

	return patch.Validate()
}

type PatchCategoryResponse CategoryDTOResponse

func (h *CategoriesHTTPHandler) PatchCategory(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	categoryID, err := core_http_request.GetUUIDPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get categoryID path value")
		return
	}

	r.Body = http.MaxBytesReader(rw, r.Body, 16*1024)

	var request PatchCategoryRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")
		return
	}

	categoryPatch := categoryPatchFromRequest(request)
	category, err := h.categoriesService.PatchCategory(ctx, categoryID, categoryPatch)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to patch category")
		return
	}

	response := PatchCategoryResponse(categoryDTOFromDomain(category))
	responseHandler.JSONResponse(response, http.StatusOK)
}

func categoryPatchFromRequest(request PatchCategoryRequest) domain.CategoryPatch {
	return domain.NewCategoryPatch(request.Name.ToDomain())
}
