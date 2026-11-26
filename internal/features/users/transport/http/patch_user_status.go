package users_transport_http

import (
	"net/http"

	core_logger "github.com/Yodau1/graphics-platform/internal/core/logger"
	core_http_request "github.com/Yodau1/graphics-platform/internal/core/transport/http/request"
	core_http_response "github.com/Yodau1/graphics-platform/internal/core/transport/http/response"
)

type PatchUserStatusRequest struct {
	IsActive *bool `json:"is_active" validate:"required"`
}

type PatchUserStatusResponse UserDTOResponse

func (h *UsersHTTPHandler) PatchUserStatus(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userID, err := core_http_request.GetUUIDPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get userID path value")
		return
	}

	r.Body = http.MaxBytesReader(rw, r.Body, 16*1024)

	var request PatchUserStatusRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")
		return
	}

	user, err := h.usersService.PatchUserStatus(ctx, userID, *request.IsActive)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to patch user status")
		return
	}

	response := PatchUserStatusResponse(userDTOFromDomain(user))
	responseHandler.JSONResponse(response, http.StatusOK)
}
