//go:generate go tool go-enum --marshal --values -f ./*.go
package types

import (
	"encoding/json"
	"html/template"
	"strings"
)

// ENUM(state, code, id_token, access_token, expires_in, refresh_token, refresh_expires_in, token_type)
type TokenResponseField string

// ENUM(info, success, error, warning)
type MessageType string

// ENUM(login.tmpl, error.tmpl, info.tmpl, form-redirect.tmpl, form-post.tmpl, user.tmpl, user-edit.tmpl, logout.tmpl)
type PageID string

func ToJsDeclaration(v any, indent int) (template.JS, error) {
	var buf strings.Builder
	enc := json.NewEncoder(&buf)
	enc.SetIndent(strings.Repeat(" ", indent), "  ")
	if err := enc.Encode(v); err != nil {
		return "", err
	}
	return template.JS(buf.String()), nil
}

type Message struct {
	Type    MessageType `json:"type"`
	Summary string      `json:"summary"`
}

type URLs struct {
	BasePath     string `json:"basePath"`
	ResourcePath string `json:"resourcePath"`
}

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Picture  string `json:"picture"`
	Email    string `json:"email"`
}

// type Client struct {
// 	ClientID string `json:"clientId"`
// 	Name     string `json:"name"`
// }

type SentinelCtx struct {
	PageID   PageID    `json:"pageId"`
	Message  *Message  `json:"message"`
	Messages []Message `json:"messages"`
	URLs     URLs      `json:"urls"`
	User     *User     `json:"user"`
	// Client   Client    `json:"client"`
}

func NewSentinelCtx(basePath string, message *Message, messages []Message, user *User) SentinelCtx {
	return SentinelCtx{
		Message:  message,
		Messages: messages,
		URLs: URLs{
			BasePath:     basePath,
			ResourcePath: basePath + "/assets",
		},
		User: user,
	}
}

type Provider struct {
	LoginURL    string `json:"loginUrl"`
	Alias       string `json:"alias"`
	ProviderID  string `json:"providerId"`
	DisplayName string `json:"displayName"`
	IconPath    string `json:"icon"`
	IsLinked    bool   `json:"isLinked"`
}

type CSRF struct {
	FieldName string `json:"fieldName"`
	Value     string `json:"value"`
}

type LoginSentinelCtx struct {
	SentinelCtx
	Providers []Provider `json:"providers"`
	CSRF      CSRF       `json:"csrf"`
}

func NewLoginSentinelCtx(sentinelCtx SentinelCtx, providers []Provider, csrf CSRF) LoginSentinelCtx {
	sentinelCtx.PageID = PageIDLogintmpl
	return LoginSentinelCtx{
		SentinelCtx: sentinelCtx,
		Providers:   providers,
		CSRF:        csrf,
	}
}

type FormRedirectSentinelCtx struct {
	SentinelCtx
	RedirectURL string `json:"redirectUrl"`
}

func NewFormRedirectSentinelCtx(sentinelCtx SentinelCtx, redirectURI string) FormRedirectSentinelCtx {
	sentinelCtx.PageID = PageIDFormRedirecttmpl
	return FormRedirectSentinelCtx{
		SentinelCtx: sentinelCtx,
		RedirectURL: redirectURI,
	}
}

type FormPostSentinelCtx struct {
	SentinelCtx
	RedirectURL string `json:"redirectUrl"`
}

func NewFormPostSentinelCtx(sentinelCtx SentinelCtx, redirectURI string) FormPostSentinelCtx {
	sentinelCtx.PageID = PageIDFormPosttmpl
	return FormPostSentinelCtx{
		SentinelCtx: sentinelCtx,
		RedirectURL: redirectURI,
	}
}

type InfoSentinelCtx struct {
	SentinelCtx
	Message Message `json:"message"`
}

func NewInfoSentinelCtx(sentinelCtx SentinelCtx, message Message) InfoSentinelCtx {
	sentinelCtx.PageID = PageIDInfotmpl
	sentinelCtx.Message = &message
	return InfoSentinelCtx{
		SentinelCtx: sentinelCtx,
		Message:     message,
	}
}

type ErrorSentinelCtx struct {
	SentinelCtx
	Message Message `json:"message"`
}

func NewErrorSentinelCtx(sentinelCtx SentinelCtx, message Message) ErrorSentinelCtx {
	sentinelCtx.PageID = PageIDErrortmpl
	sentinelCtx.Message = &message
	return ErrorSentinelCtx{
		SentinelCtx: sentinelCtx,
		Message:     message,
	}
}

type Account struct {
	Provider string `json:"provider"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Picture  string `json:"picture"`
}

type UserSentinelCtx struct {
	SentinelCtx
	User      User       `json:"user"`
	Accounts  []Account  `json:"accounts"`
	Providers []Provider `json:"providers"`
	CSRF      CSRF       `json:"csrf"`
}

func NewUserSentinelCtx(sentinelCtx SentinelCtx, user User, accounts []Account, providers []Provider, csrf CSRF) UserSentinelCtx {
	sentinelCtx.PageID = PageIDUsertmpl
	sentinelCtx.User = &user
	return UserSentinelCtx{
		SentinelCtx: sentinelCtx,
		User:        user,
		Accounts:    accounts,
		Providers:   providers,
		CSRF:        csrf,
	}
}

type UserEditSentinelCtx struct {
	SentinelCtx
	User      User       `json:"user"`
	Accounts  []Account  `json:"accounts"`
	Providers []Provider `json:"providers"`
	CSRF      CSRF       `json:"csrf"`
}

func NewUserEditSentinelCtx(sentinelCtx SentinelCtx, user User, accounts []Account, providers []Provider, csrf CSRF) UserEditSentinelCtx {
	sentinelCtx.PageID = PageIDUserEdittmpl
	sentinelCtx.User = &user
	return UserEditSentinelCtx{
		SentinelCtx: sentinelCtx,
		User:        user,
		Accounts:    accounts,
		Providers:   providers,
		CSRF:        csrf,
	}
}

type LogoutSentinelCtx struct {
	SentinelCtx
	User      User   `json:"user"`
	CSRF      CSRF   `json:"csrf"`
	Redirect  string `json:"redirect"`
	SessionID string `json:"sessionId"`
}

func NewLogoutSentinelCtx(sentinelCtx SentinelCtx, user User, redirect string, sessionID string, csrf CSRF) LogoutSentinelCtx {
	sentinelCtx.PageID = PageIDLogouttmpl
	sentinelCtx.User = &user
	return LogoutSentinelCtx{
		SentinelCtx: sentinelCtx,
		User:        user,
		Redirect:    redirect,
		SessionID:   sessionID,
		CSRF:        csrf,
	}
}
