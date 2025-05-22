package openid

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/a-h/templ"
	"github.com/akatranlp/sentinel/account"
	"github.com/akatranlp/sentinel/jose"
	"github.com/akatranlp/sentinel/openid/web/components"
	"golang.org/x/oauth2"
)

func (ip *IdentitiyProvider) ProviderCallback(w http.ResponseWriter, r *http.Request) {
	provider := r.PathValue("provider")

	p, ok := ip.providers[provider]
	if !ok {
		// TODO: Render nice frontend error-page
		http.Error(w, fmt.Sprintf("The provider %s is not configured", provider), http.StatusNotFound)
		return
	}

	ctx := r.Context()
	j := jose.GetJose(ctx)

	verifier, ok := ip.sessionManager.GetVerifier(ctx)
	if !ok {
		// TODO: Render nice frontend error-page
		http.Error(w, "invalid callback", http.StatusNotFound)
		return
	}

	oauthConfig := p.GetOauthConfig(j.Issuer() + "/" + provider + "/callback")
	authCodeOption := oauth2.VerifierOption(verifier.Verifier)

	oauthToken, err := oauthConfig.Exchange(r.Context(), r.FormValue("code"), authCodeOption)
	if err != nil {
		// TODO: Render nice frontend error-page
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var state oauthState
	if err := json.Unmarshal([]byte(r.FormValue("state")), &state); err != nil {
		// TODO: Render nice frontend error-page
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if state.RandomState != verifier.RandomState {
		// TODO: Render nice frontend error-page
		http.Error(w, "invalid state received", http.StatusInternalServerError)
		return
	}

	if err := p.ValidateToken(r.Context(), oauthToken, verifier.Nonce); err != nil {
		// TODO: Render nice frontend error-page
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	userInfo, err := p.GetUserInfo(r.Context(), oauthConfig.TokenSource(r.Context(), oauthToken))
	if err != nil {
		// TODO: Render nice frontend error-page
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var refreshExpiry time.Time
	if refreshExpiresIn, ok := oauthToken.Extra("refresh_expires_in").(int64); ok {
		refreshExpiry = time.Now().Add(time.Duration(refreshExpiresIn) * time.Second)
	}
	if refreshExpiresIn, ok := oauthToken.Extra("refresh_token_expires_in").(int64); ok {
		refreshExpiry = time.Now().Add(time.Duration(refreshExpiresIn) * time.Second)
	}

	var idToken string
	if idTokenStr, ok := oauthToken.Extra("id_token").(string); ok {
		idToken = idTokenStr
	}

	userInfo.AccountID.Provider = provider
	userInfo.AccessToken = oauthToken.AccessToken
	userInfo.Expiry = oauthToken.Expiry
	userInfo.RefreshToken = oauthToken.RefreshToken
	userInfo.RefreshExpiry = refreshExpiry
	userInfo.IDToken = idToken
	userInfo.TokenType = oauthToken.TokenType

	var user account.User
	var flashMessages []templ.Component
	if ip.sessionManager.IsAuthed(ctx) {
		userID := account.UserID(ip.sessionManager.GetAuth(ctx))
		user, err = ip.userStore.GetUserByID(r.Context(), userID)
		if err != nil {
			// TODO: Render nice frontend error-page
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if err := ip.userStore.LinkAccount(r.Context(), user.UserID, *userInfo); err != nil {
			// TODO: Render nice frontend error-page
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		accounts, err := ip.userStore.GetAccountsForUserID(r.Context(), user.UserID)
		if err != nil {
			// TODO: Render nice frontend error-page
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		for _, acc := range accounts {
			if acc.Provider == provider {
				continue
			}
			provider := ip.providers[acc.Provider]
			token, err := provider.RefreshToken(r.Context(), &oauth2.Token{AccessToken: acc.AccessToken, RefreshToken: acc.RefreshToken})
			if err != nil {
				// TODO: Add Message to manually relink account
				flashMessages = append(flashMessages, components.GitHubIcon(""))
				continue
			}

			var refreshExpiry time.Time
			if refreshExpiresIn, ok := token.Extra("refresh_expires_in").(int64); ok {
				refreshExpiry = time.Now().Add(time.Duration(refreshExpiresIn) * time.Second)
			}
			if refreshExpiresIn, ok := token.Extra("refresh_token_expires_in").(int64); ok {
				refreshExpiry = time.Now().Add(time.Duration(refreshExpiresIn) * time.Second)
			}

			var idToken string
			if idTokenStr, ok := token.Extra("id_token").(string); ok {
				idToken = idTokenStr
			}

			acc.AccessToken = token.AccessToken
			acc.Expiry = token.Expiry
			acc.RefreshToken = token.RefreshToken
			acc.RefreshExpiry = refreshExpiry
			acc.IDToken = idToken
			acc.TokenType = token.TokenType
			if err := ip.userStore.UpdateAccount(r.Context(), acc.AccountID, acc); err != nil {
				// TODO: Add Message to manually relink account
				flashMessages = append(flashMessages, components.GitHubIcon(""))
				continue
			}
		}
	} else {
		// check if user already exists
		// save tokens into database
		// and userdata if user doesnt exist yet
		acc, err := ip.userStore.GetAccountByID(r.Context(), userInfo.AccountID)
		if err == nil {
			if err = p.RevokeRefreshToken(r.Context(), acc.RefreshToken); err != nil {
				// TODO: Render nice frontend error-page
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else if !errors.Is(err, account.ErrAccountNotFound) {
			fmt.Println(err)
			// TODO: Render nice frontend error-page
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		user, err = ip.userStore.GetOrCreateUserFromAccount(r.Context(), *userInfo)
		if err != nil {
			// TODO: Render nice frontend error-page
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		accounts, err := ip.userStore.GetAccountsForUserID(r.Context(), user.UserID)
		if err != nil {
			// TODO: Render nice frontend error-page
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		for _, acc := range accounts {
			if acc.Provider == provider {
				continue
			}
			provider := ip.providers[acc.Provider]
			token, err := provider.RefreshToken(r.Context(), &oauth2.Token{AccessToken: acc.AccessToken, RefreshToken: acc.RefreshToken})
			if err != nil {
				// TODO: Add Message to manually relink account
				flashMessages = append(flashMessages, components.GitHubIcon(""))
				continue
			}

			var refreshExpiry time.Time
			if refreshExpiresIn, ok := token.Extra("refresh_expires_in").(int64); ok {
				refreshExpiry = time.Now().Add(time.Duration(refreshExpiresIn) * time.Second)
			}
			if refreshExpiresIn, ok := token.Extra("refresh_token_expires_in").(int64); ok {
				refreshExpiry = time.Now().Add(time.Duration(refreshExpiresIn) * time.Second)
			}

			var idToken string
			if idTokenStr, ok := token.Extra("id_token").(string); ok {
				idToken = idTokenStr
			}

			acc.AccessToken = token.AccessToken
			acc.Expiry = token.Expiry
			acc.RefreshToken = token.RefreshToken
			acc.RefreshExpiry = refreshExpiry
			acc.IDToken = idToken
			acc.TokenType = token.TokenType
			if err := ip.userStore.UpdateAccount(r.Context(), acc.AccountID, acc); err != nil {
				// TODO: Add Message to manually relink account
				flashMessages = append(flashMessages, components.GitHubIcon(""))
				continue
			}
		}

		ip.sessionManager.SetAuth(ctx, string(user.UserID))
	}

	authReq, ok := ip.sessionManager.GetAuthRequest(ctx)
	if ok {
		ip.sendAuthResponse(w, r, AuthTokenValues{
			AuthRequest: authReq,
			UserID:      string(user.UserID),
			AuthTime:    ip.sessionManager.GetAuthTime(ctx),
		}, flashMessages...)
		return
	}

	redirectURI, err := url.ParseRequestURI(state.Redirect)
	if err != nil {
		// TODO: Render nice frontend error-page
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := redirectURI.Query()
	if state.State != "" {
		query.Set("state", state.State)
	}
	redirectURI.RawQuery = query.Encode()

	http.Redirect(w, r, redirectURI.String(), http.StatusFound)
}
