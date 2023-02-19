package requests

type TokenRequest struct {
	GrantType string `json:"grant_type" validate:"required"`
	ClientID  string `json:"client_id" validate:"required"`
	Secret    string `json:"secret" validate:"required"`
	Username  string `json:"username" validate:"required"`
	Password  string `json:"password" validate:"required"`
}
