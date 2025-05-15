package types

type UserInfoResponse struct {
	Subject           string `json:"sub"`
	Email             string `json:"email"`
	EmailVerified     bool   `json:"email_verified"`
	Name              string `json:"name"`
	PreferredUsername string `json:"preferred_username"`
	Nickname          string `json:"nickname"`
	Picture           string `json:"picture"`
	Profile           string `json:"profile"`
}
