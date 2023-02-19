package clients

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"oauth-proxy/internal/oauth-proxy/entities"
	"oauth-proxy/internal/oauth-proxy/enums"
	"oauth-proxy/internal/oauth-proxy/handlers"
	"oauth-proxy/internal/oauth-proxy/infrastructure/repositories"
	"oauth-proxy/internal/oauth-proxy/requests"
	"oauth-proxy/internal/oauth-proxy/requests/validator"
	"oauth-proxy/internal/oauth-proxy/responses"
	"strings"
)

type Handler struct {
	Validator validator.PgValidatorI
	handlers.ClientAuthenticator
}

func NewClientsHandler(cr repositories.ClientRepositoryI, vr validator.PgValidatorI) Handler {
	return Handler{
		ClientAuthenticator: handlers.ClientAuthenticator{
			ClientRepository: cr,
		},
		Validator: vr,
	}
}

func (ch Handler) Create(ctx echo.Context) (err error) {
	clientID, secret, err := ch.ExtractBasicAuthCredentials(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusForbidden, `invalid authorization header`)
	}

	if !ch.AuthenticateClient(clientID, secret, enums.ClientManagementScope) {
		return echo.NewHTTPError(http.StatusForbidden, `forbidden, out of scope`)
	}

	var req requests.ClientCreateRequest
	if err := ctx.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if rErr := ch.Validator.Validate(req); rErr != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, rErr)
	}

	req.Name = strings.ToLower(req.Name)
	cli := entities.CreateClient(req.Name, []string{string(enums.UserManagementScope)})
	err = ch.ClientRepository.Persist(cli)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(201, responses.ClientResponse{
		ID:     cli.ID,
		Name:   cli.Name,
		Secret: cli.Secret,
		Scopes: strings.Join(cli.Scopes, " "),
	})
}
