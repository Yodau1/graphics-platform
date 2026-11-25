package resources_transport_http

import (
	"fmt"
	"net/http"
	"time"

	core_logger "github.com/Yodau1/graphics-platform/internal/core/logger"
	core_http_request "github.com/Yodau1/graphics-platform/internal/core/transport/http/request"
	core_http_response "github.com/Yodau1/graphics-platform/internal/core/transport/http/response"
)

func (h *ResourcesHTTPHandler) GetPublishedPreview(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)

	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	resourceID, err := core_http_request.GetUUIDPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get resourceID path value")
		return
	}

	fileID, err := core_http_request.GetUUIDPathValue(r, "fileId")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get fileID path value")
		return
	}

	file, content, err := h.resourcesService.GetPublishedPreview(ctx, resourceID, fileID)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get published preview")
		return
	}
	defer content.Close()

	// Выбираем тип ответа по формату, который сервер определил при загрузке
	var contentType string
	switch file.Format {
	case "png":
		contentType = "image/png"
	case "jpeg":
		contentType = "image/jpeg"
	case "glb":
		contentType = "model/gltf-binary"
	default:
		err := fmt.Errorf("unsupported preview format: %q", file.Format)
		responseHandler.ErrorResponse(err, "failed to serve published preview")
		return
	}

	rw.Header().Set("Content-Type", contentType)
	//nosniff - браузер не будет пытаться угадать тип, а будет учитывать contentType
	rw.Header().Set("X-Content-Type-Options", "nosniff")
	//просит браузер не сохранять ответ в кэше
	rw.Header().Set("Cache-Control", "no-store")

	// ServeContent читает файл и отправляет байты, в том числе по запросам Range.
	// Нулевое время отключает Last-Modified: при замене файла его CreatedAt сохраняется
	http.ServeContent(rw, r, file.OriginalName, time.Time{}, content)
}
