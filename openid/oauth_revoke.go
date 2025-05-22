package openid

import (
	"maps"
	"net/http"

	"github.com/akatranlp/sentinel/jose"
	"github.com/akatranlp/sentinel/openid/enums"
	"github.com/akatranlp/sentinel/openid/types"
	"github.com/go-viper/mapstructure/v2"
	"github.com/lestrrat-go/jwx/v3/jwt"
)

func (ip *IdentitiyProvider) OauthRevoke(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		ip.handleTokenError(w, r, &TokenError{
			ErrorType:        TokenErrorTypeInvalidRequest,
			ErrorDescription: err.Error(),
		})
		return
	}

	// Some oauth clients are sending the client id and secret as basic auth and some as form-values, we allow both
	clientID := r.FormValue(RevokeFormValueClientId.String())
	clientSecret := r.FormValue(RevokeFormValueClientSecret.String())
	if clientID == "" && clientSecret == "" {
		clientID, clientSecret, _ = r.BasicAuth()
	}
	r.Form.Set(RevokeFormValueClientId.String(), clientID)
	r.Form.Set(RevokeFormValueClientSecret.String(), clientSecret)

	params, err := mapFormValueKeys(r.Form, maps.Keys(_RevokeFormValueValue))
	if err != nil {
		ip.handleTokenError(w, r, &TokenError{
			ErrorType:        TokenErrorTypeInvalidRequest,
			ErrorDescription: err.Error(),
		})
		return
	}
	var formValues types.RevokeRequest
	if err = mapstructure.Decode(params, &formValues); err != nil {
		ip.handleTokenError(w, r, &TokenError{
			ErrorType:        TokenErrorTypeInvalidRequest,
			ErrorDescription: err.Error(),
		})
		return
	}

	reg, ok := ip.clients[formValues.ClientID]
	if !ok {
		ip.handleTokenError(w, r, &TokenError{
			ErrorType: TokenErrorTypeInvalidClient,
		})
		return
	}
	if reg.ClientSecret != "" && formValues.ClientSecret != reg.ClientSecret {
		ip.handleTokenError(w, r, &TokenError{
			ErrorType: TokenErrorTypeInvalidClient,
		})
		return
	}

	if formValues.Token == "" {
		ip.handleTokenError(w, r, &TokenError{
			ErrorType:        TokenErrorTypeInvalidRequest,
			ErrorDescription: "no token provided",
		})
		return
	}

	j := jose.GetJose(r.Context())

	token, err := jwt.ParseString(
		formValues.Token,
		jwt.WithKeySet(j.PublicKeys()),
	)
	if err != nil {
		w.WriteHeader(http.StatusOK)
		return
	}

	var tokenType string
	token.Get(enums.ClaimTokenType.String(), &tokenType)

	if tokenType == enums.OauthTokenTypeAccessToken.String() {
		ip.handleTokenError(w, r, &TokenError{
			ErrorType:        TokenErrorType("unsupported_token_type"),
			ErrorDescription: "access token are not supported",
		})
		return
	} else if tokenType == enums.OauthTokenTypeRefreshToken.String() {
		var sessionID string
		token.Get(enums.ClaimSid.String(), &sessionID)

		ip.tokenStore.RevokeSession(r.Context(), sessionID)

		w.WriteHeader(http.StatusOK)
		return
	}

	ip.handleTokenError(w, r, &TokenError{
		ErrorType:        TokenErrorTypeInvalidRequest,
		ErrorDescription: "invalid token type",
	})
	return
}

// ENUM(
// client_id,
// client_secret,
// token,
// token_hint,
// )
type RevokeFormValue string
