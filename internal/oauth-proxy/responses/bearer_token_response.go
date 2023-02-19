package responses

import (
	"encoding/json"
	"oauth-proxy/internal/oauth-proxy/entities"
	"time"
)

type BearerTokenResponse struct {
	AccessToken  entities.AccessToken
	RefreshToken entities.RefreshToken
}

func (btr BearerTokenResponse) GenerateHTTPResponse() []byte {
	respB := HTTPTokenResponse{
		AccessToken:  btr.AccessToken.ID,
		RefreshToken: btr.RefreshToken.ID,
		TokenType:    `Bearer`,
		ExpiresIn:    int(time.Since(btr.AccessToken.ExpiresAt).Seconds()),
	}

	resp, _ := json.Marshal(respB)

	return resp
}
