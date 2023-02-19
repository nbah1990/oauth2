package responses

type ClientResponse struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Secret string `json:"secret"`
	Scopes string `json:"scopes"`
}
