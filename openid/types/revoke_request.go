package types

type RevokeRequest struct {
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
	Token        string `mapstructure:"token"`
	TokenHint    string `mapstructure:"token_hint"`
}
