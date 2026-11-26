package users_transport_http

import (
	"net/http"
	"strconv"

	core_errors "github.com/Yodau1/graphics-platform/internal/core/errors"
	core_logger "github.com/Yodau1/graphics-platform/internal/core/logger"
	core_http_response "github.com/Yodau1/graphics-platform/internal/core/transport/http/response"
)

type GetUsersResponse []UserDTOResponse

func (h *UsersHTTPHandler) GetUsers(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	params := r.URL.Query()
	limit, offset := 20, 0

	if value := params.Get("limit"); value != "" {
		parsedLimit, err := strconv.Atoi(value)
		if err != nil {
			responseHandler.ErrorResponse(core_errors.ErrInvalidArgument, "limit must be an integer")
			return
		}
		limit = parsedLimit
	}

	if value := params.Get("offset"); value != "" {
		parsedOffset, err := strconv.Atoi(value)
		if err != nil {
			responseHandler.ErrorResponse(core_errors.ErrInvalidArgument, "offset must be an integer")
			return
		}
		offset = parsedOffset
	}

	users, err := h.usersService.GetUsers(ctx, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get users")
		return
	}

	response := make(GetUsersResponse, len(users))
	for i, user := range users {
		response[i] = userDTOFromDomain(user)
	}

	responseHandler.JSONResponse(response, http.StatusOK)
}
