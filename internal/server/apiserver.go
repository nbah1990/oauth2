package server

import (
	"database/sql"
	"github.com/labstack/echo/v4"
	"oauth-proxy/internal/oauth-proxy/config"
	"oauth-proxy/internal/oauth-proxy/handlers/clients"
	"oauth-proxy/internal/oauth-proxy/handlers/oauth"
	"oauth-proxy/internal/oauth-proxy/handlers/users"
	"oauth-proxy/internal/oauth-proxy/infrastructure/repositories"
	"oauth-proxy/internal/oauth-proxy/requests/validator"
)

type APIServer struct {
	Config   *config.Config
	Database *sql.DB
}

func (as APIServer) Init() (err error) {
	e := echo.New()

	as.defineRoutes(e)

	e.Logger.Fatal(e.Start(as.Config.ApplicationAddress))

	return
}

func (as APIServer) defineRoutes(e *echo.Echo) {
	vr := *validator.NewValidator()
	cr := repositories.SQLClientRepository{DB: as.Database}
	ur := repositories.SQLUserRepository{DB: as.Database}
	atr := repositories.SQLAccessTokenRepository{DB: as.Database}
	rtr := repositories.SQLRefreshTokenRepository{DB: as.Database}

	oH := oauth.NewOauthHandler(cr, ur, atr, rtr, as.Config, &vr)
	ogr := e.Group(`/oauth`)
	ogr.POST(`/token`, oH.Token)
	ogr.POST(`/token/refresh`, oH.RefreshToken)
	ogr.Any(`/token/validate`, oH.ValidateToken)

	cH := clients.NewClientsHandler(cr, vr)
	cgr := e.Group(`/clients`)
	cgr.POST(``, cH.Create)

	uH := users.NewUsersHandler(ur, cr, &vr)
	ugr := e.Group(`/users`)
	ugr.POST(``, uH.Create)
	ugr.PATCH(`/:id`, uH.Update)
}
