package requests

type UserUpdateRequest struct {
	Id       string `json:"id" validate:"required,max=100,min=3"`
	Password string `json:"password" validate:"required,max=100,min=3"`
}
