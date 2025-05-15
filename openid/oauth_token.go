package openid

import (
	"context"
	"encoding/json"
	"fmt"
	"maps"
	"net/http"
	"slices"
	"time"

	"github.com/akatranlp/go-pkg/its"
	"github.com/akatranlp/sentinel/account"
	"github.com/akatranlp/sentinel/jose"
	"github.com/akatranlp/sentinel/openid/enums"
	"github.com/akatranlp/sentinel/openid/types"
	"github.com/akatranlp/sentinel/utils"
	"github.com/go-viper/mapstructure/v2"
	"github.com/google/uuid"

	"github.com/lestrrat-go/jwx/v3/jwt"
	"golang.org/x/oauth2"
	"golang.org/x/sync/singleflight"
)

func (ip *identitiyProvider) OauthToken(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		ip.handleTokenError(w, r, &TokenError{
			ErrorType:        TokenErrorTypeInvalidRequest,
			ErrorDescription: err.Error(),
		})
		return
	}

	// Some oauth clients are sending the client id and secret as basic auth and some as form-values, we allow both
	clientID := r.FormValue(TokenFormValueClientId.String())
	clientSecret := r.FormValue(TokenFormValueClientSecret.String())
	if clientID == "" && clientSecret == "" {
		clientID, clientSecret, _ = r.BasicAuth()
	}
	r.Form.Set(TokenFormValueClientId.String(), clientID)
	r.Form.Set(TokenFormValueClientSecret.String(), clientSecret)

	var formValues types.TokenRequest
	config := utils.CreateDecoderConfig()
	config.Result = &formValues
	decoder, err := mapstructure.NewDecoder(&config)
	if err != nil {
		ip.handleTokenError(w, r, &TokenError{
			ErrorType:        TokenErrorTypeInvalidRequest,
			ErrorDescription: err.Error(),
		})
		return
	}

	params, err := mapFormValueKeys(r.Form, maps.Keys(_TokenFormValueValue))
	if err != nil {
		ip.handleTokenError(w, r, &TokenError{
			ErrorType:        TokenErrorTypeInvalidRequest,
			ErrorDescription: err.Error(),
		})
		return
	}

	if err := decoder.Decode(params); err != nil {
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

	switch formValues.GrantType {
	case enums.GrantTypeAuthorizationCode:
		authReq, ok := ip.validateCode(formValues.Code)
		if !ok {
			ip.handleTokenError(w, r, &TokenError{
				ErrorType:        TokenErrorTypeInvalidGrant,
				ErrorDescription: "invalid or expired code",
			})
			return
		}
		res, err := ip.handleAuthorizationCode(r.Context(), authReq, formValues)
		if err != nil {
			ip.handleTokenError(w, r, err)
			return
		}

		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
		return

	case enums.GrantTypeRefreshToken:
		res, err := ip.handleRefreshToken(r.Context(), formValues)
		if err != nil {
			ip.handleTokenError(w, r, err)
			return
		}

		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
		return

	case enums.GrantTypeTokenExchange:
		if reg.TokenExchangeSecret != "" && formValues.TokenExchangeSecret == reg.TokenExchangeSecret || reg.TokenExchangeSecret == "" && reg.ClientSecret != "" {
			res, err := ip.handleTokenExchange(r.Context(), formValues)
			if err != nil {
				ip.handleTokenError(w, r, err)
				return
			}

			w.Header().Set("Cache-Control", "no-store")
			w.Header().Set("Pragma", "no-cache")
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(res)
			return
		}
		ip.handleTokenError(w, r, &TokenError{
			ErrorType: TokenErrorTypeInvalidClient,
		})
		return

	default:
		ip.handleTokenError(w, r, &TokenError{
			ErrorType:        TokenErrorTypeInvalidGrant,
			ErrorDescription: "Grant-type " + formValues.GrantType.String() + " is invalid",
		})
		return
	}

}

func (ip *identitiyProvider) handleAuthorizationCode(ctx context.Context, authReq AuthTokenValues, tokenReq types.TokenRequest) (types.TokenResponse, *TokenError) {
	rURIStr := authReq.RedirectURI.String()
	if authReq.ClientID != tokenReq.ClientID {
		return types.TokenResponse{}, &TokenError{
			ErrorType:        TokenErrorTypeInvalidGrant,
			ErrorDescription: "clientID does not match with initial auth request",
		}
	}
	if rURIStr != tokenReq.RedirectURI {
		return types.TokenResponse{}, &TokenError{
			ErrorType:        TokenErrorTypeInvalidGrant,
			ErrorDescription: "redirect uri does not match with initial auth request",
		}
	}

	// Mitigate PKCE downgrade attack
	if tokenReq.CodeVerifier != "" && authReq.CodeChallengeMethod == "" {
		return types.TokenResponse{}, &TokenError{
			ErrorType:        TokenErrorTypeInvalidGrant,
			ErrorDescription: "invalid code challenge verifier",
		}
	}

	if authReq.CodeChallengeMethod == enums.CodeChallengeMethodPlain && tokenReq.CodeVerifier != authReq.CodeChallenge {
		return types.TokenResponse{}, &TokenError{
			ErrorType:        TokenErrorTypeInvalidGrant,
			ErrorDescription: "invalid code challenge verifier",
		}
	}

	if authReq.CodeChallengeMethod == enums.CodeChallengeMethodS256 &&
		oauth2.S256ChallengeFromVerifier(tokenReq.CodeVerifier) != authReq.CodeChallenge {
		return types.TokenResponse{}, &TokenError{
			ErrorType:        TokenErrorTypeInvalidGrant,
			ErrorDescription: "invalid code challenge verifier",
		}
	}

	sessionID := uuid.NewString()
	user, err := ip.userStore.GetUserByID(ctx, account.UserID(authReq.UserID))
	if err != nil {
		return types.TokenResponse{}, &TokenError{
			ErrorType:        TokenErrorTypeServerError,
			ErrorDescription: err.Error(),
		}
	}

	openID := slices.Contains(authReq.Scope, enums.ScopeOpenid)
	offlineAccess := slices.Contains(authReq.Scope, enums.ScopeOfflineAccess)

	j := jose.GetJose(ctx)

	currTime := time.Now()

	tokens, err := j.CreateTokens(jose.IDTokenCreateArg{
		TokenCreateArg: jose.TokenCreateArg{
			CurrTime:  currTime,
			Subject:   authReq.UserID,
			Audience:  []string{authReq.ClientID},
			SessionID: sessionID,
			Scope:     authReq.Scope,
		},
		Nonce:         authReq.Nonce,
		AuthTime:      authReq.AuthTime,
		Email:         user.Email,
		EmailVerified: user.EmailVerified,
		Name:          user.Name,
		Username:      user.Username,
		Picture:       user.Picture,
	}, offlineAccess, openID)

	if err != nil {
		return types.TokenResponse{}, &TokenError{
			ErrorType:        TokenErrorTypeServerError,
			ErrorDescription: err.Error(),
		}
	}

	accessExpiresIn := int(utils.Bang(tokens.AccessToken.Expiration()).Sub(currTime) / time.Second)

	var refreshExpiresIn int
	if tokens.RefreshToken != nil {
		refreshExpiresIn = int(utils.Bang(tokens.RefreshToken.Expiration()).Sub(currTime) / time.Second)

		ip.tokenStore.SetSession(ctx, sessionID, utils.Bang(tokens.RefreshToken.JwtID()), time.Now().Add(time.Duration(refreshExpiresIn)*time.Second))
	}

	return types.TokenResponse{
		AccessToken:      tokens.SignedAccessToken,
		ExpiresIn:        accessExpiresIn,
		TokenType:        "Bearer",
		IDToken:          tokens.SignedIDToken,
		RefreshToken:     tokens.SignedRefreshToken,
		RefreshExpiresIn: refreshExpiresIn,
	}, nil
}

func (ip *identitiyProvider) handleTokenError(w http.ResponseWriter, r *http.Request, err *TokenError) {
	jose := jose.GetJose(r.Context())

	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Content-Type", "application/json")

	switch err.ErrorType {
	case TokenErrorTypeInvalidRequest, TokenErrorTypeUnsupportedGrantType, TokenErrorTypeInvalidScope:
		w.WriteHeader(http.StatusBadRequest)
	case TokenErrorTypeInvalidClient:
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set("WWW-Authenticate", `Basic realm="`+jose.Issuer()+`, charset="UTF-8"`)
		return
	case TokenErrorTypeInvalidGrant, TokenErrorTypeUnauthorizedClient:
		w.WriteHeader(http.StatusUnauthorized)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(err)
}

func (ip *identitiyProvider) handleRefreshToken(ctx context.Context, req types.TokenRequest) (types.TokenResponse, *TokenError) {
	j := jose.GetJose(ctx)
	parsedRefreshToken, err := jwt.ParseString(
		req.RefreshToken,
		jwt.WithKeySet(j.PublicKeys()),
		jwt.WithIssuer(j.Issuer()),
		jwt.WithAudience(req.ClientID),
		jwt.WithClaimValue(enums.ClaimTokenType.String(), enums.OauthTokenTypeRefreshToken.String()),
	)
	if err != nil {
		return types.TokenResponse{}, &TokenError{
			ErrorType:        TokenErrorTypeInvalidGrant,
			ErrorDescription: err.Error(),
		}
	}

	var scopeStr string
	parsedRefreshToken.Get(enums.ClaimScope.String(), &scopeStr)

	var scopes enums.Scopes
	scopes.UnmarshalText([]byte(scopeStr))

	if !its.All(slices.Values(req.Scope), func(scope enums.Scope) bool { return slices.Contains(scopes, scope) }) {
		return types.TokenResponse{}, &TokenError{
			ErrorType: TokenErrorTypeInvalidScope,
		}
	}

	var sessionID string
	parsedRefreshToken.Get(enums.ClaimSid.String(), &sessionID)

	sess, err := ip.tokenStore.GetSession(ctx, sessionID)
	if err != nil {
		return types.TokenResponse{}, &TokenError{
			ErrorType:        TokenErrorTypeInvalidGrant,
			ErrorDescription: err.Error(),
		}
	}
	if sess.RefreshJTI != utils.Bang(parsedRefreshToken.JwtID()) {
		ip.tokenStore.RevokeSession(ctx, sessionID)
		return types.TokenResponse{}, &TokenError{
			ErrorType:        TokenErrorTypeInvalidGrant,
			ErrorDescription: "token is revoked",
		}
	}

	userID := utils.Bang(parsedRefreshToken.Subject())

	// TODO: Go oder nicht go das ist hier die Frage

	// go ip.refreshAccountToken(account.UserID(userID))
	ip.refreshAccountToken(account.UserID(userID))

	currTime := time.Now()
	arg := jose.TokenCreateArg{
		CurrTime:  currTime,
		Subject:   userID,
		Audience:  utils.Bang(parsedRefreshToken.Audience()),
		SessionID: sessionID,
		Scope:     scopes,
	}

	accessToken, signedAccessToken, err := j.CreateAccessToken(arg)
	if err != nil {
		return types.TokenResponse{}, &TokenError{
			ErrorType:        TokenErrorTypeServerError,
			ErrorDescription: err.Error(),
		}
	}

	refreshToken, signedRefreshToken, err := j.CreateRefreshToken(arg)
	if err != nil {
		return types.TokenResponse{}, &TokenError{
			ErrorType:        TokenErrorTypeServerError,
			ErrorDescription: err.Error(),
		}
	}

	accessExpiresIn := int(utils.Bang(accessToken.Expiration()).Sub(currTime) / time.Second)
	refreshExpiresIn := int(utils.Bang(refreshToken.Expiration()).Sub(currTime) / time.Second)
	tokenID := utils.Bang(refreshToken.JwtID())

	if err := ip.tokenStore.SetSession(ctx, sessionID, tokenID, time.Now().Add(time.Duration(refreshExpiresIn)*time.Second)); err != nil {
		return types.TokenResponse{}, &TokenError{
			ErrorType:        TokenErrorTypeServerError,
			ErrorDescription: err.Error(),
		}
	}

	return types.TokenResponse{
		AccessToken:      signedAccessToken,
		ExpiresIn:        accessExpiresIn,
		TokenType:        "Bearer",
		RefreshToken:     signedRefreshToken,
		RefreshExpiresIn: refreshExpiresIn,
	}, nil
}

func (ip *identitiyProvider) handleTokenExchange(ctx context.Context, req types.TokenRequest) (types.TokenResponse, *TokenError) {
	if req.RequestedTokenType != enums.OauthTokenTypeAccessToken || req.SubjectTokenType != enums.OauthTokenTypeAccessToken {
		return types.TokenResponse{}, &TokenError{
			ErrorType:        TokenErrorTypeInvalidGrant,
			ErrorDescription: enums.ErrInvalidOauthTokenType.Error(),
		}
	}
	j := jose.GetJose(ctx)

	parsedRefreshToken, err := jwt.ParseString(
		req.SubjectToken,
		jwt.WithKeySet(j.PublicKeys()),
		jwt.WithIssuer(j.Issuer()),
		jwt.WithAudience(req.ClientID),
		jwt.WithClaimValue(enums.ClaimTokenType.String(), enums.OauthTokenTypeAccessToken.String()),
	)
	if err != nil {
		return types.TokenResponse{}, &TokenError{
			ErrorType:        TokenErrorTypeInvalidGrant,
			ErrorDescription: err.Error(),
		}
	}

	var sessionID string
	parsedRefreshToken.Get(enums.ClaimSid.String(), &sessionID)

	if _, err = ip.tokenStore.GetSession(ctx, sessionID); err != nil {
		return types.TokenResponse{}, &TokenError{
			ErrorType:        TokenErrorTypeInvalidGrant,
			ErrorDescription: err.Error(),
		}
	}

	userID := utils.Bang(parsedRefreshToken.Subject())

	acc, err := ip.userStore.GetAccountByProvider(ctx, account.UserID(userID), req.RequestedIssuer)
	if err != nil {
		return types.TokenResponse{}, &TokenError{
			ErrorType:        TokenErrorTypeInvalidTarget,
			ErrorDescription: err.Error(),
		}
	}

	return types.TokenResponse{
		AccessToken:     acc.AccessToken,
		ExpiresIn:       int(acc.Expiry.Sub(time.Now()) / time.Second),
		TokenType:       acc.TokenType,
		IssuedTokenType: req.RequestedTokenType,
	}, nil
}

var g singleflight.Group

func (ip *identitiyProvider) refreshAccountToken(userID account.UserID) {
	_, err, _ := g.Do(string(userID), func() (any, error) {
		ctx := context.Background()
		accounts, err := ip.userStore.GetAccountsForUserID(ctx, userID)
		if err != nil {
			return nil, err
		}

		for _, acc := range accounts {
			p, ok := ip.providers[acc.Provider]
			if !ok {
				fmt.Println("TODO: could this happen?")
				return nil, err
			}
			token, err := p.RefreshToken(ctx, &oauth2.Token{
				AccessToken:  acc.AccessToken,
				RefreshToken: acc.RefreshToken,
				Expiry:       acc.Expiry,
				TokenType:    acc.TokenType,
			})
			if err != nil {
				fmt.Println("Refresh Error", acc.Provider, err)
				return nil, err
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
			acc.TokenType = token.TokenType
			if idToken != "" {
				acc.IDToken = idToken
			}

			if err := ip.userStore.UpdateAccount(ctx, acc.AccountID, acc); err != nil {
				fmt.Println(err)
				return nil, err
			}
		}
		return nil, nil
	})
	if err != nil {
		return
	}
}

// ENUM(
// client_id,
// client_secret,
// redirect_uri,
// grant_type,
// code,
// code_verifier,
// scope,
// refresh_token,
// requested_token_type,
// requested_issuer,
// subject_token,
// subject_token_type,
// __token_exchange_secret__,
// )
type TokenFormValue string

// ENUM(
// invalid_request,
// invalid_client,
// invalid_grant,
// unauthorized_client,
// unsupported_grant_type,
// invalid_scope,
// server_error,
// invalid_target,
// )
type TokenErrorType string

type TokenError struct {
	ErrorType        TokenErrorType `json:"error"`
	ErrorDescription string         `json:"error_description,omitzero"`
}
