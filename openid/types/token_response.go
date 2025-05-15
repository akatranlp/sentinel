package types

import "github.com/akatranlp/sentinel/openid/enums"

type TokenResponse struct {
	AccessToken      string               `json:"access_token"`
	ExpiresIn        int                  `json:"expires_in"`
	RefreshToken     string               `json:"refresh_token,omitzero"`
	RefreshExpiresIn int                  `json:"refresh_expires_in,omitzero"`
	IDToken          string               `json:"id_token,omitzero"`
	TokenType        string               `json:"token_type"`
	IssuedTokenType  enums.OauthTokenType `json:"issued_token_type"`
}
