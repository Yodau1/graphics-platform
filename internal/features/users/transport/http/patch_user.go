package users_transport_http

import (
	"net/http"

	"github.com/Yodau1/graphics-platform/internal/core/domain"
	core_logger "github.com/Yodau1/graphics-platform/internal/core/logger"
	core_http_request "github.com/Yodau1/graphics-platform/internal/core/transport/http/request"
	core_http_response "github.com/Yodau1/graphics-platform/internal/core/transport/http/response"
	core_http_types "github.com/Yodau1/graphics-platform/internal/core/transport/http/types"
)

type PatchUserRequest struct {
	Name  core_http_types.Nullable[string] `json:"name" swaggertype:"string"`
	Email core_http_types.Nullable[string] `json:"email" swaggertype:"string"`
}

func (r *PatchUserRequest) Validate() error {

	patch := userPatchFromRequest(*r)

	return patch.Validate()
}

type PatchUserResponse UserDTOResponse

func (h *UsersHTTPHandler) PatchUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userID, err := core_http_request.GetUUIDPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get userID path value",
		)

		return
	}

	r.Body = http.MaxBytesReader(rw, r.Body, 16*1024)

	var request PatchUserRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate HTTP request",
		)

		return
	}

	userPatch := userPatchFromRequest(request)

	userDomain, err := h.usersService.PatchUser(ctx, userID, userPatch)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to patch user",
		)

		return
	}

	response := PatchUserResponse(userDTOFromDomain(userDomain))

	responseHandler.JSONResponse(response, http.StatusOK)
}

// userPatchFromRequest преобразует HTTP-запрос в доменный патч,
// сохраняя различие между отсутствующим полем, null и переданным значением.
func userPatchFromRequest(request PatchUserRequest) domain.UserPatch {
	return domain.NewUserPatch(
		request.Name.ToDomain(),
		request.Email.ToDomain(),
	)
}
