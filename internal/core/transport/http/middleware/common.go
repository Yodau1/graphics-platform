package core_http_middleware

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	core_logger "github.com/Yodau1/graphics-platform/internal/core/logger"
	core_http_response "github.com/Yodau1/graphics-platform/internal/core/transport/http/response"
)

const requestIDHeader = "X-Request-ID"

// RequestID — middleware, обеспечивающий каждый запрос уникальным идентификатором.
// Если клиент передаёт X-Request-ID — используем его (полезно для распределённой трассировки).
// Иначе генерируем новый X-Request-ID.
// Идентификатор добавляется и в заголовок ответа, чтобы клиент мог его использовать.
func RequestID() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(requestIDHeader)
			if requestID == "" {
				requestID = uuid.NewString()
			}

			r.Header.Set(requestIDHeader, requestID)
			w.Header().Set(requestIDHeader, requestID)

			next.ServeHTTP(w, r)
		})
	}
}

// Logger — middleware, кладущий логгер в контекст запроса.
// Обогащает логгер полями request_id и url, чтобы все последующие
// обработчики автоматически логировали эти поля.
func Logger(log *core_logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(requestIDHeader)
			l := log.With(
				zap.String("request_id", requestID),
				zap.String("url", r.URL.String()),
			)

			ctx := core_logger.ToContext(r.Context(), l)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Trace — middleware для логирования входящих запросов и времени их обработки.
// Использует ResponseWriter-обёртку, чтобы перехватить статус-код ответа,
// который иначе недоступен после вызова WriteHeader.
func Trace() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := core_logger.FromContext(r.Context())
			rw := core_http_response.NewResponseWriter(w)

			before := time.Now()
			log.Debug(
				">>> incoming HTTP request",
				zap.String("http_method", r.Method),
				zap.Time("time", before.UTC()),
			)

			next.ServeHTTP(rw, r)

			log.Debug(
				"<<< done HTTP request",
				zap.Int("status_code", rw.GetStatusCode()),
				zap.Duration("latency", time.Since(before)),
			)
		})
	}
}

// Panic — middleware для перехвата паник и возврата HTTP 500.
// Без этого middleware паника в обработчике уронила бы всю горутину,
// а стандартная библиотека Go вернула бы пустой ответ клиенту.
func Panic() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := core_logger.FromContext(r.Context())
			responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

			defer func() {
				if p := recover(); p != nil {
					responseHandler.PanicResponse(
						p,
						"during handle HTTP request got unexpected panic",
					)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
