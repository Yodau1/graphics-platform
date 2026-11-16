package core_http_middleware

import (
	"net/http"
)

//Middleware оборачивает HTTP-обработчик.
// ChainMiddleware собирает обёртки так,
// чтобы первая в списке выполнялась первой.
// Logger добавляет к логгеру метод запроса и передаёт
// его дальше через r.Context()

type Middleware func(http.Handler) http.Handler

func ChainMiddleware(handler http.Handler, middleware ...Middleware) http.Handler {
	for i := len(middleware) - 1; i >= 0; i-- {
		handler = middleware[i](handler)
	}

	return handler
}
