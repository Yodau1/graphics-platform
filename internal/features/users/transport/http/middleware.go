package users_transport_http

import (
	"net/http"
	"strings"

	"github.com/Yodau1/graphics-platform/internal/core/domain"
	core_errors "github.com/Yodau1/graphics-platform/internal/core/errors"
	core_logger "github.com/Yodau1/graphics-platform/internal/core/logger"
	core_http_middleware "github.com/Yodau1/graphics-platform/internal/core/transport/http/middleware"
	core_http_request "github.com/Yodau1/graphics-platform/internal/core/transport/http/request"
	core_http_response "github.com/Yodau1/graphics-platform/internal/core/transport/http/response"
)

// authorizeUser создает middleware для проверки доступа к операциям у пользователей
// adminOnly = true — доступ разрешён только администратору
// adminOnly = false — администратору разрешён доступ к любому профилю,
// а обычному пользователю и автору — только к своему
func (h *UsersHTTPHandler) authorizeUser(adminOnly bool) core_http_middleware.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := core_logger.FromContext(ctx)
			responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

			// Разбиваем значение заголовка по пробельным символам
			//Нужно, чтобы достать тип токена и сам токен
			parts := strings.Fields(r.Header.Get("Authorization"))

			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				responseHandler.ErrorResponse(
					core_errors.ErrUnauthorized,
					"bearer token required",
				)
				return
			}

			//из сервиса получаем пользователя по токену сессии
			user, err := h.usersService.Authenticate(ctx, parts[1])
			if err != nil {
				responseHandler.ErrorResponse(err, "failed to authenticate user")
				return
			}

			if user.Role == domain.UserRoleAdmin {
				next.ServeHTTP(rw, r)
				return
			}

			if adminOnly {
				responseHandler.ErrorResponse(
					core_errors.ErrForbidden,
					"administrator role required",
				)
				return
			}

			// Для пользователя и автора читаем ID профиля из URL. /users/{id}
			userID, err := core_http_request.GetUUIDPathValue(r, "id")
			if err != nil {
				responseHandler.ErrorResponse(
					err,
					"failed to get userID path value",
				)
				return
			}

			// Сравниваем владельца токена с владельцем запрошенного профиля
			if user.ID != userID {
				responseHandler.ErrorResponse(
					core_errors.ErrForbidden,
					"profile belongs to another user",
				)
				return
			}

			//Если проверки пройдены, передаем запрос обработчику
			next.ServeHTTP(rw, r)
		})
	}
}
