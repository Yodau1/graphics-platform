package resources_transport_http

import (
	"net/http"

	core_errors "github.com/Yodau1/graphics-platform/internal/core/errors"
	core_logger "github.com/Yodau1/graphics-platform/internal/core/logger"
	core_http_middleware "github.com/Yodau1/graphics-platform/internal/core/transport/http/middleware"
	core_http_request "github.com/Yodau1/graphics-platform/internal/core/transport/http/request"
	core_http_response "github.com/Yodau1/graphics-platform/internal/core/transport/http/response"
)

type GetResourceFilesResponse []ResourceFileDTOResponse

// GetResourceFiles возвращает действующие файлы указанной карточки.
func (h *ResourcesHTTPHandler) GetResourceFiles(rw http.ResponseWriter, r *http.Request) {
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

	files, err := h.resourcesService.GetResourceFiles(ctx, user, resourceID)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get resource files")
		return
	}

	// Собираем ответ без внутреннего ключа хранения.
	response := make(GetResourceFilesResponse, len(files))
	for i, file := range files {
		response[i] = resourceFileDTOFromDomain(file)
	}

	responseHandler.JSONResponse(response, http.StatusOK)
}
