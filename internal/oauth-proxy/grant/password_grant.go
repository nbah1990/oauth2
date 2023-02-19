package grant

import (
	"errors"
	"oauth-proxy/internal/oauth-proxy/infrastructure/repositories"
	"oauth-proxy/internal/oauth-proxy/requests"
	"oauth-proxy/internal/oauth-proxy/responses"
)

type PasswordGrant struct {
	Grant
	ClientRepository repositories.ClientRepositoryI
	UserRepository   repositories.UserRepositoryI
}

func (pg PasswordGrant) AccessTokenRequest(req requests.TokenRequest) (res responses.AuthorizationResponse, err error) {
	c, err := pg.ClientRepository.GetByID(req.ClientID)
	if err != nil {
		return nil, errors.New(`client not found`)
	}

	if !c.CheckSecret(req.Secret) || c.Revoked {
		return nil, errors.New(`client not found`)
	}

	u, err := pg.UserRepository.FindByUsername(req.Username)
	if err != nil {
		return nil, errors.New(`user not found`)
	}

	if !u.CheckPassword(req.Password) {
		return nil, errors.New(`user not found`)
	}

	at, err := pg.IssueAccessToken(c, u)
	if err != nil {
		return nil, errors.New(`at internal error`)
	}

	rt, err := pg.IssueRefreshToken(at)
	if err != nil {
		return nil, errors.New(`rt internal error`)
	}

	res = responses.HTTPTokenResponse{
		AccessToken:  at.ID,
		RefreshToken: rt.ID,
		TokenType:    `Bearer`,
		ExpiresIn:    int(pg.GetConfig().AccessTokenExpirationPeriod.Seconds()),
	}

	return res, nil
}
