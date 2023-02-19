package grant

import (
	"errors"
	"oauth-proxy/internal/oauth-proxy/infrastructure/repositories"
	"oauth-proxy/internal/oauth-proxy/requests"
	"oauth-proxy/internal/oauth-proxy/responses"
)

type RefreshTokenGrant struct {
	Grant
	AccessTokenRepository  repositories.AccessTokenRepositoryI
	RefreshTokenRepository repositories.RefreshTokenRepositoryI

	ClientRepository repositories.ClientRepositoryI
	UserRepository   repositories.UserRepositoryI
}

func (rtg RefreshTokenGrant) RefreshTokenRequest(req requests.RefreshTokenRequest) (res responses.AuthorizationResponse, err error) {
	c, err := rtg.ClientRepository.GetByID(req.ClientID)
	if err != nil {
		return nil, errors.New(`client not found`)
	}

	if !c.CheckSecret(req.Secret) || c.Revoked {
		return nil, errors.New(`client not found`)
	}

	oldRt, err := rtg.RefreshTokenRepository.GetByID(req.RefreshToken)

	if err != nil {
		return nil, errors.New(`refresh token not found`)
	}

	if oldRt.IsExpired() || oldRt.Revoked {
		return nil, errors.New(`refresh token expired`)
	}

	oldAt, err := rtg.AccessTokenRepository.GetByID(oldRt.AccessTokenID)
	if err != nil {
		return nil, errors.New(`access token not found`)
	}

	u, err := rtg.UserRepository.GetByID(oldAt.UserID)
	if err != nil {
		return nil, errors.New(`user not found`)
	}

	at, err := rtg.IssueAccessToken(c, u)
	if err != nil {
		return nil, errors.New(`at internal error`)
	}

	rt, err := rtg.IssueRefreshToken(at)
	if err != nil {
		return nil, errors.New(`rt internal error`)
	}

	_ = rtg.AccessTokenRepository.Revoke(oldAt)
	_ = rtg.RefreshTokenRepository.Revoke(oldRt)

	res = responses.HTTPTokenResponse{
		AccessToken:  at.ID,
		RefreshToken: rt.ID,
		TokenType:    `Bearer`,
		ExpiresIn:    int(rtg.GetConfig().AccessTokenExpirationPeriod.Seconds()),
	}

	return res, nil
}
