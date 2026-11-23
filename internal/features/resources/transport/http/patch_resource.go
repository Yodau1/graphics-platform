package resources_transport_http

import (
	"net/http"

	"github.com/google/uuid"

	"github.com/Yodau1/graphics-platform/internal/core/domain"
	core_errors "github.com/Yodau1/graphics-platform/internal/core/errors"
	core_logger "github.com/Yodau1/graphics-platform/internal/core/logger"
	core_http_middleware "github.com/Yodau1/graphics-platform/internal/core/transport/http/middleware"
	core_http_request "github.com/Yodau1/graphics-platform/internal/core/transport/http/request"
	core_http_response "github.com/Yodau1/graphics-platform/internal/core/transport/http/response"
	core_http_types "github.com/Yodau1/graphics-platform/internal/core/transport/http/types"
)

type PatchResourceRequest struct {
	CategoryID  core_http_types.Nullable[uuid.UUID] `json:"category_id" swaggertype:"string"`
	Title       core_http_types.Nullable[string]    `json:"title" swaggertype:"string"`
	Description core_http_types.Nullable[string]    `json:"description" swaggertype:"string"`
}

// Validate проверяет, что патч содержит допустимый набор полей.
func (r *PatchResourceRequest) Validate() error {
	patch := resourcePatchFromRequest(*r)

	return patch.Validate()
}

type PatchResourceResponse ResourceDTOResponse

func (h *ResourcesHTTPHandler) PatchResource(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

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

	r.Body = http.MaxBytesReader(rw, r.Body, 16*1024)

	var request PatchResourceRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")
		return
	}

	// Сервис получает сам патч и пользователя из проверенной сессии.
	resourcePatch := resourcePatchFromRequest(request)
	resource, err := h.resourcesService.PatchResource(ctx, user, resourceID, resourcePatch)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to patch resource")
		return
	}

	response := PatchResourceResponse(resourceDTOFromDomain(resource))
	responseHandler.JSONResponse(response, http.StatusOK)
}

func resourcePatchFromRequest(request PatchResourceRequest) domain.ResourcePatch {
	return domain.NewResourcePatch(
		request.CategoryID.ToDomain(),
		request.Title.ToDomain(),
		request.Description.ToDomain(),
	)
}
