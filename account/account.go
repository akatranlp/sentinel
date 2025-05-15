package account

import "time"

type UserID string

type User struct {
	UserID
	Name          string
	Username      string
	Picture       string
	Email         string
	EmailVerified bool
}

type AccountID struct {
	Provider   string
	ProviderID string
}

type Account struct {
	AccountID
	Email             string
	EmailVerified     bool
	AccessToken       string
	Expiry            time.Time
	RefreshToken      string
	RefreshExpiry     time.Time
	TokenType         string
	IDToken           string
	Name              string
	PreferredUsername string
	Nickname          string
	Picture           string
	Profile           string
	UserID            UserID
}
