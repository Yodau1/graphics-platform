package users_transport_http

import (
	"net/http"
	"time"

	core_logger "github.com/Yodau1/graphics-platform/internal/core/logger"
	core_http_request "github.com/Yodau1/graphics-platform/internal/core/transport/http/request"
	core_http_response "github.com/Yodau1/graphics-platform/internal/core/transport/http/response"
)

// структура для входящего http запроса
type LoginRequest struct {
	Email    string `json:"email" validate:"required,max=254"`
	Password string `json:"password" validate:"required"`
}

// структура для http ответа
type LoginResponse struct {
	AccessToken string    `json:"access_token"`
	TokenType   string    `json:"token_type"`
	ExpiresAt   time.Time `json:"expires_at"`
}

func (h *UsersHTTPHandler) Login(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	r.Body = http.MaxBytesReader(rw, r.Body, 16*1024)

	var request LoginRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")
		return
	}

	token, expiresAt, err := h.usersService.Login(ctx, request.Email, request.Password)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to sign up")
		return
	}

	//Заприщаем кэширование ответа, т.к он содержит токен доступа
	rw.Header().Set("Cache-Control", "no-store")

	//Bearer-токен. Сервер будет доверять любому, кто предоставляет данный вид токена
	response := LoginResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresAt:   expiresAt,
	}

	responseHandler.JSONResponse(response, http.StatusOK)

}
