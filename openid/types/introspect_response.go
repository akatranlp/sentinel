package types

import "github.com/akatranlp/sentinel/openid/enums"

type IntrospectResponse struct {
	Active     bool         `json:"active"`
	Scope      enums.Scopes `json:"scope,omitzero"`
	ClientID   string       `json:"client_id,omitzero"`
	Username   string       `json:"username,omitzero"`
	TokenType  string       `json:"token_type,omitzero"`
	Expiration int64        `json:"exp,omitzero"`
	IssuedAt   int64        `json:"iat,omitzero"`
	NotBefore  int64        `json:"nbf,omitzero"`
	Subject    string       `json:"sub,omitzero"`
	Audience   []string     `json:"aud,omitzero"`
	Issuer     string       `json:"iss,omitzero"`
	Jti        string       `json:"jti,omitzero"`
}
