package types

import "github.com/akatranlp/sentinel/openid/enums"

type OpenIDConfiguration struct {
	Issuer                string `json:"issuer"`
	AuthorizationEndpoint string `json:"authorization_endpoint"`
	TokenEndpoint         string `json:"token_endpoint"`
	RevocationEndpoint    string `json:"revocation_endpoint"`
	IntrospectionEndpoint string `json:"introspection_endpoint"`
	UserInfoEndpoint      string `json:"userinfo_endpoint"`
	JWKSURI               string `json:"jwks_uri"`
	EndSessionEndpoint    string `json:"end_session_endpoint"`

	ScopesSupported                           []enums.Scope                  `json:"scopes_supported"`
	ResponseTypesSupported                    []enums.ResponseType           `json:"response_types_supported"`
	ResponseModesSupported                    []enums.ResponseMode           `json:"response_modes_supported"`
	GrantTypesSupported                       []enums.GrantType              `json:"grant_types_supported"`
	TokenEndpointAuthMethodsSupported         []enums.EndpointAuthMethod     `json:"token_endpoint_auth_methods_supported"`
	IntrospectionEndpointAuthMethodsSupported []enums.EndpointAuthMethod     `json:"introspection_endpoint_auth_methods_supported"`
	RevocationEndpointAuthMethodsSupported    []enums.EndpointAuthMethod     `json:"revocation_endpoint_auth_methods_supported"`
	SubjectTypesSupported                     []enums.SubjectType            `json:"subject_types_supported"`
	IDTokenSigningAlgValuesSupported          []enums.IDTokenSigningAlgValue `json:"id_token_signing_alg_values_supported"`
	ClaimTypesSupported                       []enums.ClaimType              `json:"claim_types_supported"`
	ClaimsSupported                           []enums.Claim                  `json:"claims_supported"`
	CodeChallengeMethodsSupported             []enums.CodeChallengeMethod    `json:"code_challenge_methods_supported"`
}

func CreateOIDCConfig(origin string) OpenIDConfiguration {
	return OpenIDConfiguration{
		Issuer:                origin,
		AuthorizationEndpoint: origin + "/oauth/authorize",
		TokenEndpoint:         origin + "/oauth/token",
		RevocationEndpoint:    origin + "/oauth/revoke",
		IntrospectionEndpoint: origin + "/oauth/introspect",
		UserInfoEndpoint:      origin + "/oauth/userinfo",
		JWKSURI:               origin + "/oauth/discovery/keys",
		EndSessionEndpoint:    origin + "/oauth/logout",

		ScopesSupported:                           enums.ScopeValues(),
		ResponseTypesSupported:                    enums.ResponseTypeValues(),
		ResponseModesSupported:                    enums.ResponseModeValues(),
		GrantTypesSupported:                       enums.GrantTypeValues(),
		TokenEndpointAuthMethodsSupported:         enums.EndpointAuthMethodValues(),
		IntrospectionEndpointAuthMethodsSupported: enums.EndpointAuthMethodValues(),
		RevocationEndpointAuthMethodsSupported:    enums.EndpointAuthMethodValues(),
		SubjectTypesSupported:                     enums.SubjectTypeValues(),
		IDTokenSigningAlgValuesSupported:          enums.IDTokenSigningAlgValueValues(),
		ClaimTypesSupported:                       enums.ClaimTypeValues(),
		ClaimsSupported:                           enums.ClaimValues(),
		CodeChallengeMethodsSupported:             enums.CodeChallengeMethodValues(),
	}
}
