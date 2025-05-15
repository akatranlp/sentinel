package types

import "github.com/akatranlp/sentinel/openid/enums"

type TokenRequest struct {
	ClientID     string          `mapstructure:"client_id"`
	ClientSecret string          `mapstructure:"client_secret"`
	GrantType    enums.GrantType `mapstructure:"grant_type"`
	// Autorization Code GrantType
	RedirectURI  string       `mapstructure:"redirect_uri"`
	Code         string       `mapstructure:"code"`
	CodeVerifier string       `mapstructure:"code_verifier"`
	Scope        enums.Scopes `mapstructure:"scope"`
	// Refresh Token GrantType
	RefreshToken string `mapstructure:"refresh_token"`
	// Token Exchange GrantType
	RequestedTokenType  enums.OauthTokenType `mapstructure:"requested_token_type"`
	RequestedIssuer     string               `mapstructure:"requested_issuer"`
	SubjectToken        string               `mapstructure:"subject_token"`
	SubjectTokenType    enums.OauthTokenType `mapstructure:"subject_token_type"`
	TokenExchangeSecret string               `mapstructure:"__token_exchange_secret__"`
}
