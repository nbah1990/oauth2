package requests

type ClientCreateRequest struct {
	Name string `json:"name" validate:"required,alpha,max=20,min=3"`
}
