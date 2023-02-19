package responses

import "encoding/json"

type AuthorizationResponse interface {
	GenerateHTTPResponse() []byte
}

type HTTPTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

func (htr HTTPTokenResponse) GenerateHTTPResponse() []byte {
	res, err := json.Marshal(htr)
	if err != nil {
		return nil
	}

	return res
}
