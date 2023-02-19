package responses

type UserResponse struct {
	ID         string `json:"id"`
	Username   string `json:"username"`
	ExternalID string `json:"external_id"`
}
