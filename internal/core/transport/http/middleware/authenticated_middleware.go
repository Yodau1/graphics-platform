package core_http_middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/Yodau1/graphics-platform/internal/core/domain"
	core_errors "github.com/Yodau1/graphics-platform/internal/core/errors"
	core_logger "github.com/Yodau1/graphics-platform/internal/core/logger"
	core_http_response "github.com/Yodau1/graphics-platform/internal/core/transport/http/response"
)

type SessionAuthenticator interface {
	Authenticate(ctx context.Context, token string) (domain.User, error)
}

// authenticatedUserKey — закрытый ключ для хранения пользователя в контексте.
// Собственный тип ключа исключает совпадения с ключами других пакетов
type authenticatedUserKey struct{}

// Middleware для аутентификации
func RequireAuthenticatedUser(authenticator SessionAuthenticator) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := core_logger.FromContext(ctx)
			responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

			// Заголовок должен содержать схему Bearer и сам токен
			parts := strings.Fields(r.Header.Get("Authorization"))
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				responseHandler.ErrorResponse(core_errors.ErrUnauthorized, "bearer token required")
				return
			}

			user, err := authenticator.Authenticate(ctx, parts[1])
			if err != nil {
				responseHandler.ErrorResponse(err, "failed to authenticate user")
				return
			}

			// Передаём обработчику пользователя из проверенной сессии
			ctx = context.WithValue(ctx, authenticatedUserKey{}, user)
			next.ServeHTTP(rw, r.WithContext(ctx))
		})
	}
}

// Проверка на авторизованного пользователя
func AuthenticatedUserFromContext(ctx context.Context) (domain.User, bool) {
	user, ok := ctx.Value(authenticatedUserKey{}).(domain.User)

	return user, ok
}
