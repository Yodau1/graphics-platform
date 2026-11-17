package users_transport_http

import (
	"net/http"

	core_logger "github.com/Yodau1/graphics-platform/internal/core/logger"
	core_http_request "github.com/Yodau1/graphics-platform/internal/core/transport/http/request"
	core_http_response "github.com/Yodau1/graphics-platform/internal/core/transport/http/response"
)

type CreateUserRequest struct {
	Name     string `json:"name"     validate:"required,min=3,max=100"`
	Email    string `json:"email"    validate:"required,max=254,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type CreateUserResponse UserDTOResponse

func (h *UsersHTTPHandler) CreateUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	// Ограничиваем размер JSON-запроса регистрации до 16 КиБ.
	r.Body = http.MaxBytesReader(rw, r.Body, 16*1024)

	var request CreateUserRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")
		return
	}

	userDomain, err := h.usersService.CreateUser(
		ctx,
		request.Name,
		request.Email,
		request.Password,
	)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create user")
		return
	}

	response := CreateUserResponse(userDTOFromDomain(userDomain))
	responseHandler.JSONResponse(response, http.StatusCreated)
}
