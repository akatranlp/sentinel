//go:generate go tool go-enum --marshal
package provider

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/akatranlp/sentinel/account"
	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

// ENUM(gitlab, github, gitea)
type ProviderType string

type Provider interface {
	GetOauthConfig(redirectURL string) oauth2.Config
	GetUserInfo(ctx context.Context, tokenSource oauth2.TokenSource) (*account.Account, error)
	ValidateToken(ctx context.Context, token *oauth2.Token, nonce string) error
	RevokeRefreshToken(ctx context.Context, token string) error
	RefreshToken(ctx context.Context, token *oauth2.Token) (*oauth2.Token, error)
	GetType() ProviderType

	GetName() string
	GetSlug() string
}

var _ Provider = (*OIDCProvider)(nil)
var _ Provider = (*OauthProvider)(nil)

type OIDCProvider struct {
	Name         string
	Slug         string
	Type         ProviderType
	ClientID     string
	ClientSecret string
	Provider     *oidc.Provider
	Verifier     *oidc.IDTokenVerifier
	Scopes       []string
}

func NewGitLabProvider(name, slug, clientID, clientSecret, baseURL string) (*OIDCProvider, error) {
	return NewOIDCProvider(name, slug, ProviderTypeGitlab, clientID, clientSecret, baseURL, []string{oidc.ScopeOpenID, "profile", "email", "api"})
}

func NewGiteaProvider(name, slug, clientID, clientSecret, baseURL string) (*OIDCProvider, error) {
	return NewOIDCProvider(name, slug, ProviderTypeGitea, clientID, clientSecret, baseURL, []string{oidc.ScopeOpenID, "profile", "email"})
}

func NewOIDCProvider(name, slug string, providerType ProviderType, clientID, clientSecret, baseURL string, scopes []string) (*OIDCProvider, error) {
	provider, err := oidc.NewProvider(context.TODO(), baseURL)
	if err != nil {
		return nil, err
	}
	internalScopes := make([]string, len(scopes))
	copy(internalScopes, scopes)

	if !slices.Contains(internalScopes, oidc.ScopeOpenID) {
		internalScopes = append(internalScopes, oidc.ScopeOpenID)
	}
	return &OIDCProvider{
		Name:         name,
		Slug:         slug,
		Type:         providerType,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Provider:     provider,
		Verifier:     provider.Verifier(&oidc.Config{ClientID: clientID}),
		Scopes:       internalScopes,
	}, nil
}

func (p *OIDCProvider) GetName() string       { return p.Name }
func (p *OIDCProvider) GetSlug() string       { return p.Slug }
func (p *OIDCProvider) GetType() ProviderType { return p.Type }

func (p *OIDCProvider) RevokeRefreshToken(ctx context.Context, token string) error {
	type providerClaims struct {
		RevokeEndpoint string `json:"revocation_endpoint"`
	}
	var claims providerClaims
	if err := p.Provider.Claims(&claims); err != nil {
		return err
	}

	if claims.RevokeEndpoint == "" {
		return nil
	}

	params := make(url.Values)
	params.Set("token", token)
	params.Set("token_hint", "refresh_token")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, claims.RevokeEndpoint, strings.NewReader(params.Encode()))
	if err != nil {
		return err
	}
	req.SetBasicAuth(p.ClientID, p.ClientSecret)
	res, err := oauth2.NewClient(ctx, nil).Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		data, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		return errors.New(string(data))
	}

	return nil
}

func (p *OIDCProvider) GetOauthConfig(redirectURL string) oauth2.Config {
	return oauth2.Config{
		ClientID:     p.ClientID,
		ClientSecret: p.ClientSecret,
		RedirectURL:  redirectURL,
		Endpoint:     p.Provider.Endpoint(),
		Scopes:       p.Scopes,
	}
}

func (p *OIDCProvider) RefreshToken(ctx context.Context, token *oauth2.Token) (*oauth2.Token, error) {
	oauthConfig := p.GetOauthConfig("")
	token.Expiry = time.Now().Add(-1 * time.Hour)

	token, err := oauthConfig.TokenSource(ctx, token).Token()
	if err != nil {
		return nil, err
	}

	if rawIDToken, ok := token.Extra("id_token").(string); ok {
		if _, err := p.Verifier.Verify(ctx, rawIDToken); err != nil {
			return nil, err
		}
	}

	return token, nil
}

type UserInfoData struct {
	Subject           string `json:"sub"`
	Name              string `json:"name"`
	Nickname          string `json:"nickname"`
	PreferredUsername string `json:"preferred_username"`
	Email             string `json:"email"`
	EmailVerified     bool   `json:"email_verified"`
	Profile           string `json:"profile"`
	Picture           string `json:"picture"`
}

func (p *OIDCProvider) GetUserInfo(ctx context.Context, tokenSource oauth2.TokenSource) (*account.Account, error) {
	userInfo, err := p.Provider.UserInfo(ctx, tokenSource)
	if err != nil {
		return nil, err
	}

	var data UserInfoData
	if err := userInfo.Claims(&data); err != nil {
		return nil, err
	}

	return &account.Account{
		AccountID:         account.AccountID{ProviderID: data.Subject},
		Email:             data.Email,
		EmailVerified:     data.EmailVerified,
		Name:              data.Name,
		Nickname:          data.Nickname,
		PreferredUsername: data.PreferredUsername,
		Profile:           data.Profile,
		Picture:           data.Picture,
	}, nil
}

func (p *OIDCProvider) ValidateToken(ctx context.Context, token *oauth2.Token, nonce string) error {
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return errors.New("empty token")
	}

	idToken, err := p.Verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return err
	}

	if idToken.Nonce != nonce {
		return errors.New("nonce verification is invalid")
	}

	if idToken.AccessTokenHash != "" {
		if err := idToken.VerifyAccessToken(token.AccessToken); err != nil {
			return err
		}
	}

	return nil
}

type OauthProvider struct {
	Name           string
	Slug           string
	Type           ProviderType
	ClientID       string
	ClientSecret   string
	Endpoint       oauth2.Endpoint
	UserInfoGetter UserInfoGetter
	Scopes         []string
}

type UserInfoGetter = func(ctx context.Context, token oauth2.TokenSource) (*account.Account, error)

func githubUserInfoGetter(ctx context.Context, tokenSource oauth2.TokenSource) (*account.Account, error) {
	client := oauth2.NewClient(ctx, tokenSource)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.github.com/user", nil)
	if err != nil {
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		data, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(string(data))
	}

	type githubUser struct {
		ID        int    `json:"id"`
		Name      string `json:"name"`
		HTMLURL   string `json:"html_url"`
		AvatarURL string `json:"avatar_url"`
		Login     string `json:"login"`
	}

	var ghUser githubUser

	if err := json.NewDecoder(res.Body).Decode(&ghUser); err != nil {
		return nil, err
	}

	req, err = http.NewRequest(http.MethodGet, "https://api.github.com/user/emails", nil)
	if err != nil {
		return nil, err
	}
	res, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		data, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(string(data))
	}

	type githubEmail struct {
		Email   string `json:"email"`
		Primary bool   `json:"primary"`
	}

	var ghEmails []githubEmail

	if err := json.NewDecoder(res.Body).Decode(&ghEmails); err != nil {
		return nil, err
	}

	emailIdx := slices.IndexFunc(ghEmails, func(v githubEmail) bool { return v.Primary })
	return &account.Account{
		AccountID:         account.AccountID{ProviderID: strconv.Itoa(ghUser.ID)},
		Email:             ghEmails[emailIdx].Email,
		EmailVerified:     true,
		Profile:           ghUser.HTMLURL,
		Name:              ghUser.Name,
		PreferredUsername: ghUser.Login,
		Nickname:          ghUser.Login,
		Picture:           ghUser.AvatarURL,
	}, nil
}

func NewGitHubProvider(name, slug, clientID, clientSecret string) (*OauthProvider, error) {
	return NewOauthProvider(name, slug, ProviderTypeGithub, clientID, clientSecret, github.Endpoint, []string{"profile", "email"}, githubUserInfoGetter)
}

func NewOauthProvider(name, slug string, providerType ProviderType, clientID, clientSecret string, endpoint oauth2.Endpoint, scopes []string, userInfoGetter UserInfoGetter) (*OauthProvider, error) {
	internalScopes := make([]string, len(scopes))
	copy(internalScopes, scopes)

	return &OauthProvider{
		Name:           name,
		Slug:           slug,
		Type:           providerType,
		ClientID:       clientID,
		ClientSecret:   clientSecret,
		Endpoint:       endpoint,
		UserInfoGetter: userInfoGetter,
		Scopes:         internalScopes,
	}, nil
}

func (p *OauthProvider) GetName() string       { return p.Name }
func (p *OauthProvider) GetSlug() string       { return p.Slug }
func (p *OauthProvider) GetType() ProviderType { return p.Type }

func (p *OauthProvider) RevokeRefreshToken(ctx context.Context, token string) error { return nil }

func (p *OauthProvider) RefreshToken(ctx context.Context, token *oauth2.Token) (*oauth2.Token, error) {
	oauthConfig := p.GetOauthConfig("")
	token.Expiry = time.Now().Add(-1 * time.Hour)
	return oauthConfig.TokenSource(ctx, token).Token()
}

func (p *OauthProvider) GetOauthConfig(redirectURL string) oauth2.Config {
	return oauth2.Config{
		ClientID:     p.ClientID,
		ClientSecret: p.ClientSecret,
		RedirectURL:  redirectURL,
		Endpoint:     p.Endpoint,
		Scopes:       p.Scopes,
	}
}

func (p *OauthProvider) GetUserInfo(ctx context.Context, tokenSource oauth2.TokenSource) (*account.Account, error) {
	return p.UserInfoGetter(ctx, tokenSource)
}

func (p *OauthProvider) ValidateToken(ctx context.Context, token *oauth2.Token, nonce string) error {
	return nil
}

type FactoryParams struct {
	Type         string
	Name         string
	Slug         string
	BaseURL      string
	ClientID     string
	ClientSecret string
}

func ProviderFactory(params FactoryParams) (Provider, error) {
	switch params.Type {
	case "github":
		return NewGitHubProvider(params.Name, params.Slug, params.ClientID, params.ClientSecret)
	case "gitea":
		return NewGiteaProvider(params.Name, params.Slug, params.ClientID, params.ClientSecret, params.BaseURL)
	case "gitlab":
		return NewGitLabProvider(params.Name, params.Slug, params.ClientID, params.ClientSecret, params.BaseURL)
	default:
		return nil, errors.New("no valid provider type")
	}
}
