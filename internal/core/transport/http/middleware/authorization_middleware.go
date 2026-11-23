package core_http_middleware

import (
	"net/http"

	"github.com/Yodau1/graphics-platform/internal/core/domain"
	core_errors "github.com/Yodau1/graphics-platform/internal/core/errors"
	core_logger "github.com/Yodau1/graphics-platform/internal/core/logger"
	core_http_response "github.com/Yodau1/graphics-platform/internal/core/transport/http/response"
)

// Middleware на проверку администратора
func RequireAdmin() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := core_logger.FromContext(ctx)
			responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

			user, ok := AuthenticatedUserFromContext(ctx)
			if !ok {
				responseHandler.ErrorResponse(core_errors.ErrUnauthorized, "authenticated user required")
				return
			}

			if user.Role != domain.UserRoleAdmin {
				responseHandler.ErrorResponse(core_errors.ErrForbidden, "administrator role required")
				return
			}

			//передаем запрос дальше, если это админ
			next.ServeHTTP(rw, r)
		})
	}
}
