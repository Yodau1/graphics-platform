package resources_transport_http

import (
	"io"
	"net/http"

	"github.com/Yodau1/graphics-platform/internal/core/domain"
	core_errors "github.com/Yodau1/graphics-platform/internal/core/errors"
	core_logger "github.com/Yodau1/graphics-platform/internal/core/logger"
	core_http_middleware "github.com/Yodau1/graphics-platform/internal/core/transport/http/middleware"
	core_http_request "github.com/Yodau1/graphics-platform/internal/core/transport/http/request"
	core_http_response "github.com/Yodau1/graphics-platform/internal/core/transport/http/response"
	resources_service "github.com/Yodau1/graphics-platform/internal/features/resources/service"
)

type UploadResourceFileResponse ResourceFileDTOResponse

// UploadResourceFile принимает один файл из multipart-запроса.
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

	// Оставляем 64 КиБ для заголовков и границ multipart.
	r.Body = http.MaxBytesReader(rw, r.Body, resources_service.MaxResourceFileSize+64*1024)
	if err := r.ParseMultipartForm(1 << 20); err != nil {
		responseHandler.ErrorResponse(core_errors.ErrInvalidArgument, "invalid multipart form or file too large")
		return
	}
	defer r.MultipartForm.RemoveAll()

	kinds := r.MultipartForm.Value["kind"]
	files := r.MultipartForm.File["file"]
	if len(kinds) != 1 || len(files) != 1 {
		responseHandler.ErrorResponse(core_errors.ErrInvalidArgument, "one kind and one file are required")
		return
	}

	if files[0].Size > resources_service.MaxResourceFileSize {
		responseHandler.ErrorResponse(core_errors.ErrInvalidArgument, "file exceeds 10 MiB")
		return
	}

	source, err := files[0].Open()
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to open uploaded file")
		return
	}
	defer source.Close()

	// Считываем максимум на один байт больше лимита. Сервис проверит итоговый размер.
	data, err := io.ReadAll(io.LimitReader(source, resources_service.MaxResourceFileSize+1))
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to read uploaded file")
		return
	}

	file, err := h.resourcesService.UploadResourceFile(
		ctx,
		user,
		resourceID,
		domain.ResourceFileKind(kinds[0]),
		files[0].Filename,
		data,
	)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to upload resource file")
		return
	}

	response := UploadResourceFileResponse(resourceFileDTOFromDomain(file))
	responseHandler.JSONResponse(response, http.StatusCreated)
}
