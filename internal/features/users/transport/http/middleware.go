package users_transport_http

import (
	"net/http"

	"github.com/Yodau1/graphics-platform/internal/core/domain"
	core_errors "github.com/Yodau1/graphics-platform/internal/core/errors"
	core_logger "github.com/Yodau1/graphics-platform/internal/core/logger"
	core_http_middleware "github.com/Yodau1/graphics-platform/internal/core/transport/http/middleware"
	core_http_request "github.com/Yodau1/graphics-platform/internal/core/transport/http/request"
	core_http_response "github.com/Yodau1/graphics-platform/internal/core/transport/http/response"
)

// Разрешает чтение и изменение владельцу либо администратору.
func requireProfileAccess() core_http_middleware.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := core_logger.FromContext(ctx)
			responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)
			user, ok := core_http_middleware.AuthenticatedUserFromContext(ctx)
			if !ok {
				responseHandler.ErrorResponse(core_errors.ErrUnauthorized, "authenticated user required")
				return
			}

			if user.Role == domain.UserRoleAdmin {
				next.ServeHTTP(rw, r)
				return
			}

			userID, err := core_http_request.GetUUIDPathValue(r, "id")
			if err != nil {
				responseHandler.ErrorResponse(err, "failed to get userID path value")
				return
			}

			if user.ID != userID {
				// У пользователя есть действующая сессия, но профиль чужой.
				responseHandler.ErrorResponse(core_errors.ErrForbidden, "profile belongs to another user")
				return
			}

			next.ServeHTTP(rw, r)
		})
	}
}
