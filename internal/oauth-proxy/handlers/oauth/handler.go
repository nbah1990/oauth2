package oauth

import (
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
	"oauth-proxy/internal/oauth-proxy/config"
	"oauth-proxy/internal/oauth-proxy/enums"
	"oauth-proxy/internal/oauth-proxy/grant"
	"oauth-proxy/internal/oauth-proxy/infrastructure/repositories"
	"oauth-proxy/internal/oauth-proxy/requests"
	"oauth-proxy/internal/oauth-proxy/requests/validator"
	"strings"
)

type Handler struct {
	passwordGrant     grant.PasswordGrant
	refreshTokenGrant grant.RefreshTokenGrant

	accessTokenRepository repositories.AccessTokenRepositoryI
	validator             *validator.PgValidator
}

func NewOauthHandler(
	cr repositories.ClientRepositoryI,
	ur repositories.UserRepositoryI,
	atr repositories.AccessTokenRepositoryI,
	rtr repositories.RefreshTokenRepositoryI,
	cfg *config.Config,
	vr *validator.PgValidator,
) Handler {
	gr := grant.NewAbstractGrant(atr, rtr, cfg)
	return Handler{
		passwordGrant: grant.PasswordGrant{
			Grant:            gr,
			ClientRepository: cr,
			UserRepository:   ur,
		},
		refreshTokenGrant: grant.RefreshTokenGrant{
			Grant:                  gr,
			ClientRepository:       cr,
			UserRepository:         ur,
			AccessTokenRepository:  atr,
			RefreshTokenRepository: rtr,
		},
		accessTokenRepository: atr,
		validator:             vr,
	}
}

func (oh *Handler) Token(ctx echo.Context) (err error) {
	var req requests.TokenRequest
	if err := ctx.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if rErr := oh.validator.Validate(req); rErr != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, rErr)
	}

	if req.GrantType != string(enums.PasswordGrantType) {
		return ctx.JSON(http.StatusBadRequest, errors.New(`invalid grant type`))
	}

	res, err := oh.passwordGrant.AccessTokenRequest(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusForbidden, err.Error())
	}

	return ctx.JSON(200, res)
}

func (oh *Handler) RefreshToken(ctx echo.Context) (err error) {
	var req requests.RefreshTokenRequest

	if err := ctx.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if rErr := oh.validator.Validate(req); rErr != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, rErr)
	}

	if req.GrantType != string(enums.RefreshTokenGrantType) {
		return ctx.JSON(http.StatusBadRequest, errors.New(`invalid grant type`))
	}

	res, err := oh.refreshTokenGrant.RefreshTokenRequest(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusForbidden, err.Error())
	}

	return ctx.JSON(200, res)
}

func (oh *Handler) ValidateToken(ctx echo.Context) (err error) {
	token, err := ExtractBearerToken(ctx)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, errors.New(`unauthorized`).Error())
	}

	at, err := oh.accessTokenRepository.GetByID(token)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, errors.New(`unauthorized`).Error())
	}

	if at.IsExpired() {
		return ctx.JSON(http.StatusUnauthorized, errors.New(`token expired`).Error())
	}

	ctx.Response().Header().Set("X-User-Id", at.UserID)

	return ctx.JSON(200, struct {
		UserID string `json:"user_id"`
	}{UserID: at.UserID})
}

func ExtractBearerToken(ctx echo.Context) (string, error) {
	auth := ctx.Request().Header.Get("Authorization")
	if auth == "" {
		return "", errors.New("authorization header is missing")
	}

	parts := strings.Split(auth, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("invalid authorization header format")
	}

	return parts[1], nil
}
