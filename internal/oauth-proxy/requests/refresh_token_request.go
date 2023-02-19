package requests

type RefreshTokenRequest struct {
	GrantType    string `json:"grant_type" validate:"required"`
	ClientID     string `json:"client_id" validate:"required"`
	Secret       string `json:"secret" validate:"required"`
	RefreshToken string `json:"refresh_token" validate:"required"`
}
