package users_transport_http

import (
	"net/http"
	"strings"

	core_errors "github.com/Yodau1/graphics-platform/internal/core/errors"
	core_logger "github.com/Yodau1/graphics-platform/internal/core/logger"
	core_http_response "github.com/Yodau1/graphics-platform/internal/core/transport/http/response"
)

// DTO для ответа
type BecomeAuthorResponse UserDTOResponse

func (h *UsersHTTPHandler) BecomeAuthor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	parts := strings.Fields(r.Header.Get("Authorization"))
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		responseHandler.ErrorResponse(core_errors.ErrUnauthorized, "bearer token required")
		return
	}

	userDomain, err := h.usersService.BecomeAuthor(ctx, parts[1])
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to become an author")
		return
	}

	response := BecomeAuthorResponse(userDTOFromDomain(userDomain))
	responseHandler.JSONResponse(response, http.StatusOK)

}
