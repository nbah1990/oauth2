package enums

type GrantType string

const (
	PasswordGrantType     GrantType = `password`
	RefreshTokenGrantType GrantType = `refresh_token`
)
