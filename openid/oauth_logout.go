package openid

import (
	"fmt"
	"maps"
	"net/http"
	"net/url"

	"github.com/akatranlp/sentinel/jose"
	"github.com/akatranlp/sentinel/openid/enums"
	"github.com/akatranlp/sentinel/openid/types"
	"github.com/akatranlp/sentinel/openid/web"
	csrf "github.com/akatranlp/sentinel/session/gorilla_csrf"
	"github.com/akatranlp/sentinel/utils"
	"github.com/go-viper/mapstructure/v2"
	"github.com/lestrrat-go/jwx/v3/jwt"
	"github.com/lestrrat-go/jwx/v3/jwt/openid"
)

func (ip *identitiyProvider) OauthLogout(w http.ResponseWriter, r *http.Request) {
	// sess, err := c.GetAuthSession(r)
	// if err != nil {
	// 	// TODO: Send a nice frontend back
	// 	http.Error(w, "invalid session", http.StatusBadRequest)
	// 	return
	// }
	var err error

	if r.Method == http.MethodPost && r.FormValue(ip.sessionManager.CsrfFormField()) != "" {
		redirect := r.FormValue("redirect")
		sessionID := r.FormValue("sid")

		if sessionID != "" {
			ip.tokenStore.RevokeSession(r.Context(), sessionID)
		}

		if err = ip.sessionManager.Destroy(r.Context()); err != nil {
			// TODO: Better Logging
			fmt.Println("session destroy did not work", err)
		}

		if redirect == "" {
			http.Error(w, "TODO: UNIMPLEMENTED FINISHED LOGOUT PAGE", http.StatusNotImplemented)
			return
		}

		http.Redirect(w, r, redirect, http.StatusFound)
		return
	}

	if err = r.ParseForm(); err != nil {
		ip.handleLogoutError("", w, r, err.Error(), http.StatusBadRequest)
		return
	}

	var formValues types.LogoutRequest
	if err = mapstructure.Decode(mapFormValueKeysWithoutError(r.Form, maps.Keys(_LogoutFormValueValue)), &formValues); err != nil {
		ip.handleLogoutError("", w, r, err.Error(), http.StatusBadRequest)
		return
	}

	j := jose.GetJose(r.Context())

	var clientID string
	var sessionID string
	var canDirectlyLogout bool

	if formValues.IDTokenHint != "" {
		options := []jwt.ParseOption{
			jwt.WithResetValidators(true),
			jwt.WithKeySet(j.PublicKeys()),
			jwt.WithToken(openid.New()),
			jwt.WithIssuer(j.Issuer()),
		}

		if formValues.ClientID != "" {
			options = append(options, jwt.WithAudience(formValues.ClientID))
		}

		var token jwt.Token
		token, err = jwt.ParseString(formValues.IDTokenHint, options...)
		if err != nil {
			ip.handleLogoutError("", w, r, err.Error(), http.StatusBadRequest)
			return
		}

		clientID = utils.Bang(token.Audience())[0]
		token.Get(enums.ClaimSid.String(), &sessionID)
		canDirectlyLogout = true
	} else if formValues.ClientID != "" {
		clientID = formValues.ClientID
	}

	if clientID == "" {
		ip.handleLogoutError(sessionID, w, r, "no client id", http.StatusBadRequest)
		return
	}

	reg, ok := ip.clients[clientID]
	if !ok {
		ip.handleLogoutError(sessionID, w, r, "no registration", http.StatusBadRequest)
		return
	}

	var postLogoutRedirectURI *url.URL
	if formValues.PostLogoutRedirectURI != "" {
		var redirectURI *url.URL
		redirectURI, err = url.ParseRequestURI(formValues.PostLogoutRedirectURI)
		if err != nil {
			ip.handleLogoutError(sessionID, w, r, err.Error(), http.StatusBadRequest)
			return
		}

		if err = reg.CheckPostLogoutRedirectURI(*redirectURI); err != nil {
			ip.handleLogoutError(sessionID, w, r, "Check redirect uri failed: "+err.Error(), http.StatusBadRequest)
			return
		}
		postLogoutRedirectURI = redirectURI
	}

	if postLogoutRedirectURI != nil {
		if formValues.State != "" {
			params := postLogoutRedirectURI.Query()
			params.Set("state", formValues.State)
			postLogoutRedirectURI.RawQuery = params.Encode()
		}
	}

	var redirect string
	if postLogoutRedirectURI != nil {
		redirect = postLogoutRedirectURI.String()
	}

	if !ip.sessionManager.IsAuthed(r.Context()) {
		if postLogoutRedirectURI != nil && canDirectlyLogout {
			ip.tokenStore.RevokeSession(r.Context(), sessionID)
			http.Redirect(w, r, postLogoutRedirectURI.String(), http.StatusFound)
			return
		} else if canDirectlyLogout {
			ip.tokenStore.RevokeSession(r.Context(), sessionID)
			http.Error(w, "TODO: UNIMPLEMENTED FINISHED LOGOUT PAGE", http.StatusNotImplemented)
			return
		} else {
			web.Logout(ip.sessionManager.CsrfFormField(), csrf.Token(r), sessionID, redirect).Render(r.Context(), w)
			return
		}
	}

	if postLogoutRedirectURI != nil && canDirectlyLogout {
		ip.tokenStore.RevokeSession(r.Context(), sessionID)
		if err = ip.sessionManager.Destroy(r.Context()); err != nil {
			// TODO: Better Logging
			fmt.Println("session destroy did not work", err)
		}
		http.Redirect(w, r, postLogoutRedirectURI.String(), http.StatusFound)
		return
	} else if canDirectlyLogout {
		web.Logout(ip.sessionManager.CsrfFormField(), csrf.Token(r), sessionID, redirect).Render(r.Context(), w)
		return
	} else {
		// TODO: Add userinfo to logout page
		web.Logout(ip.sessionManager.CsrfFormField(), csrf.Token(r), sessionID, redirect).Render(r.Context(), w)
		return
	}

}

func (ip *identitiyProvider) handleLogoutError(sessionID string, w http.ResponseWriter, r *http.Request, errMsg string, statusCode int) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")
	w.WriteHeader(statusCode)
	csrfToken := csrf.Token(r)
	_ = csrfToken
	// redirectURI = ""
	// Add csrf
	// Add sessionID
	fmt.Println("ERROR: MESSAGE: ", errMsg)
	web.Logout(ip.sessionManager.CsrfFormField(), csrf.Token(r), sessionID, "").Render(r.Context(), w)
	return
}

// ENUM(
// client_id,
// id_token_hint,
// post_logout_redirect_uri,
// state,
// )
type LogoutFormValue string
