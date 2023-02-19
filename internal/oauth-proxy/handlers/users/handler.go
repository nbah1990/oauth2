package users

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"oauth-proxy/internal/oauth-proxy/entities"
	"oauth-proxy/internal/oauth-proxy/enums"
	"oauth-proxy/internal/oauth-proxy/handlers"
	"oauth-proxy/internal/oauth-proxy/infrastructure/repositories"
	"oauth-proxy/internal/oauth-proxy/requests"
	"oauth-proxy/internal/oauth-proxy/requests/validator"
	"oauth-proxy/internal/oauth-proxy/responses"
	"oauth-proxy/internal/oauth-proxy/services/hash"
	"strings"
)

type Handler struct {
	UserRepository repositories.UserRepositoryI
	Validator      *validator.PgValidator
	handlers.ClientAuthenticator
}

func NewUsersHandler(ur repositories.UserRepositoryI, cr repositories.ClientRepositoryI, vr *validator.PgValidator) Handler {
	return Handler{
		UserRepository: ur,
		Validator:      vr,
		ClientAuthenticator: handlers.ClientAuthenticator{
			ClientRepository: cr,
		},
	}
}

func (uh Handler) Create(ctx echo.Context) (err error) {
	var req requests.UserCreateRequest
	if err := ctx.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if rErr := uh.Validator.Validate(req); rErr != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, rErr)
	}

	clientID, secret, err := uh.ExtractBasicAuthCredentials(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusForbidden, `invalid authorization header`)
	}

	if !uh.AuthenticateClient(clientID, secret, enums.UserManagementScope) {
		return echo.NewHTTPError(http.StatusForbidden, `forbidden, out of scope`)
	}

	req.Username = strings.ToLower(req.Username)
	u := entities.CreateUser(req.Username, req.Password, req.ExternalID)
	err = uh.UserRepository.Persist(u)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(201, responses.UserResponse{
		ID:         u.ID,
		Username:   u.Username,
		ExternalID: u.ExternalID,
	})
}

func (uh Handler) Update(ctx echo.Context) (err error) {
	id := ctx.Param("id")

	var req requests.UserUpdateRequest
	if err := ctx.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	req.Id = id

	if rErr := uh.Validator.Validate(req); rErr != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, rErr)
	}

	clientID, secret, err := uh.ExtractBasicAuthCredentials(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusForbidden, `invalid authorization header`)
	}

	if !uh.AuthenticateClient(clientID, secret, enums.UserManagementScope) {
		return echo.NewHTTPError(http.StatusForbidden, `forbidden, out of scope`)
	}

	u, err := uh.UserRepository.GetByID(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, `not found`)
	}

	u.Password, err = hash.BcryptHash(req.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, `password saving error`)
	}

	err = uh.UserRepository.Persist(u)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf(`user saving error:%s`, err.Error()))
	}

	return ctx.JSON(200, responses.UserResponse{
		ID:         u.ID,
		Username:   u.Username,
		ExternalID: u.ExternalID,
	})
}
