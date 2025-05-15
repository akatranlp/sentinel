package types

type LogoutRequest struct {
	IDTokenHint           string `mapstructure:"id_token_hint"`
	PostLogoutRedirectURI string `mapstructure:"post_logout_redirect_uri"`
	ClientID              string `mapstructure:"client_id"`
	State                 string `mapstructure:"state"`
}
