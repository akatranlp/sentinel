package openid

import (
	"bytes"
	"encoding/json"
	"maps"
	"net/http"

	"github.com/akatranlp/sentinel/account"
	"github.com/akatranlp/sentinel/jose"
	"github.com/akatranlp/sentinel/openid/enums"
	"github.com/akatranlp/sentinel/openid/types"
	"github.com/akatranlp/sentinel/utils"

	"github.com/go-viper/mapstructure/v2"
	"github.com/lestrrat-go/jwx/v3/jwt"
)

func (ip *IdentitiyProvider) OauthIntrospect(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{
			"error":             "invalid_request",
			"error_description": err.Error(),
		})
		return
	}

	// Some oauth clients are sending the client id and secret as basic auth and some as form-values, we allow both
	clientID := r.FormValue(IntrospectFormValueClientId.String())
	clientSecret := r.FormValue(IntrospectFormValueClientSecret.String())
	if clientID == "" && clientSecret == "" {
		clientID, clientSecret, _ = r.BasicAuth()
	}
	r.Form.Set(IntrospectFormValueClientId.String(), clientID)
	r.Form.Set(IntrospectFormValueClientSecret.String(), clientSecret)

	var formValues types.IntrospectRequest
	if err := mapstructure.Decode(mapFormValueKeysWithoutError(r.Form, maps.Keys(_IntrospectFormValueValue)), &formValues); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{
			"error":             "invalid_request",
			"error_description": err.Error(),
		})
		return
	}

	reg, ok := ip.clients[formValues.ClientID]
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]any{
			"error":             "invalid_client",
			"error_description": "invalid client id or client secret",
		})
		return
	}
	if reg.ClientSecret != "" && formValues.ClientSecret != reg.ClientSecret {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]any{
			"error":             "invalid_client",
			"error_description": "invalid client id or client secret",
		})
		return
	}

	if formValues.Token == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{
			"error":             "invalid_request",
			"error_description": "token was empty",
		})
		return
	}

	j := jose.GetJose(r.Context())

	w.Header().Set("Content-Type", "application/json")
	var res types.IntrospectResponse

	token, err := jwt.ParseString(
		formValues.Token,
		jwt.WithKeySet(j.PublicKeys()),
		jwt.WithAudience(formValues.ClientID),
		jwt.WithIssuer(j.Issuer()),
	)
	if err != nil {
		json.NewEncoder(w).Encode(res)
		return
	}

	userID := utils.Bang(token.Subject())
	var sessionID string
	token.Get(enums.ClaimSid.String(), &sessionID)

	if _, err = ip.userStore.GetUserByID(r.Context(), account.UserID(userID)); err != nil {
		json.NewEncoder(w).Encode(res)
		return
	}

	if _, err = ip.tokenStore.GetSession(r.Context(), sessionID); err != nil {
		json.NewEncoder(w).Encode(res)
		return
	}

	var buf bytes.Buffer
	if err = json.NewEncoder(&buf).Encode(token); err != nil {
		json.NewEncoder(w).Encode(res)
		return
	}
	if err = json.NewDecoder(&buf).Decode(&res); err != nil {
		json.NewEncoder(w).Encode(res)
		return
	}
	res.Active = true
	res.ClientID = formValues.ClientID
	json.NewEncoder(w).Encode(res)
}

// ENUM(
// token,
// token_hint,
// client_id,
// client_secret,
// )
type IntrospectFormValue string
