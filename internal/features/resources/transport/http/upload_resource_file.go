package resources_transport_http

import (
	"net/http"

	core_errors "github.com/Yodau1/graphics-platform/internal/core/errors"
	core_logger "github.com/Yodau1/graphics-platform/internal/core/logger"
	core_http_middleware "github.com/Yodau1/graphics-platform/internal/core/transport/http/middleware"
	core_http_request "github.com/Yodau1/graphics-platform/internal/core/transport/http/request"
	core_http_response "github.com/Yodau1/graphics-platform/internal/core/transport/http/response"
)

type UploadResourceFileResponse ResourceFileDTOResponse

// UploadResourceFile добавляет файл в существующую карточку.
func (h *ResourcesHTTPHandler) UploadResourceFile(rw http.ResponseWriter, r *http.Request) {
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

	form, err := readResourceFileForm(rw, r, true)
	if err != nil {
		// ErrorResponse преобразует ошибку разбора формы в HTTP-ответ.
		responseHandler.ErrorResponse(err, "failed to read upload form")
		return
	}

	file, err := h.resourcesService.UploadResourceFile(
		ctx,
		user,
		resourceID,
		form.kind,
		form.originalName,
		form.data,
	)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to upload resource file")
		return
	}

	response := UploadResourceFileResponse(resourceFileDTOFromDomain(file))
	responseHandler.JSONResponse(response, http.StatusCreated)
}
