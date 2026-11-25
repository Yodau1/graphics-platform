package resources_transport_http

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"github.com/Yodau1/graphics-platform/internal/core/domain"
	core_logger "github.com/Yodau1/graphics-platform/internal/core/logger"
	core_http_request "github.com/Yodau1/graphics-platform/internal/core/transport/http/request"
	core_http_response "github.com/Yodau1/graphics-platform/internal/core/transport/http/response"
)

// PublishedPreviewDTO содержит данные, нужные интерфейсу для открытия файла.
type PublishedPreviewDTO struct {
	ID     uuid.UUID               `json:"id"`
	Kind   domain.ResourceFileKind `json:"kind"`
	Format string                  `json:"format"`
	URL    string                  `json:"url"`
}

type GetPublishedResourceResponse struct {
	ResourceDTOResponse
	Previews []PublishedPreviewDTO `json:"previews"`
}

// GetPublishedResource открывает публичную карточку со ссылками на её превью.
func (h *ResourcesHTTPHandler) GetPublishedResource(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	resourceID, err := core_http_request.GetUUIDPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get resourceID path value")
		return
	}

	// GetPublishedResource получает карточку и метаданные её превью.
	resource, previews, err := h.resourcesService.GetPublishedResource(ctx, resourceID)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get published resource")
		return
	}

	response := GetPublishedResourceResponse{
		ResourceDTOResponse: resourceDTOFromDomain(resource),
		Previews:            make([]PublishedPreviewDTO, len(previews)),
	}

	for i, preview := range previews {
		response.Previews[i] = PublishedPreviewDTO{
			ID:     preview.ID,
			Kind:   preview.Kind,
			Format: preview.Format,
			// Sprintf подставляет два UUID в путь уже существующего обработчика файла.
			URL: fmt.Sprintf("/api/v1/catalog/resources/%s/previews/%s", resource.ID, preview.ID),
		}
	}

	responseHandler.JSONResponse(response, http.StatusOK)
}
