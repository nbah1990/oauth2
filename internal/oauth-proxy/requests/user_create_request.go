package requests

type UserCreateRequest struct {
	Username   string `json:"username" validate:"required,alpha,max=20,min=3"`
	Password   string `json:"password" validate:"required,max=100,min=3"`
	ExternalID string `json:"external_id" validate:"max=100"`
}
