package users_transport_http

import (
	"net/http"
	"strings"

	core_errors "github.com/Yodau1/graphics-platform/internal/core/errors"
	core_logger "github.com/Yodau1/graphics-platform/internal/core/logger"
	core_http_response "github.com/Yodau1/graphics-platform/internal/core/transport/http/response"
)

// Logout отзывает сессию, указанную в заголовке Authorization.
func (h *UsersHTTPHandler) Logout(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	// Разбираем заголовок вида «Bearer <токен>».
	parts := strings.Fields(r.Header.Get("Authorization"))
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		responseHandler.ErrorResponse(core_errors.ErrUnauthorized, "bearer token required")
		return
	}

	// Отзываем токен из запроса. Ошибку БД передаём обработчику ответов.
	if err := h.usersService.Logout(ctx, parts[1]); err != nil {
		responseHandler.ErrorResponse(err, "failed to sign out")
		return
	}

	responseHandler.NoContentResponse()
}
