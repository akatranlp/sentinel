package openid

import (
	"maps"
	"net/http"
	"net/url"
	"time"

	"github.com/a-h/templ"
	"github.com/akatranlp/sentinel/openid/enums"
	"github.com/akatranlp/sentinel/openid/types"
	"github.com/akatranlp/sentinel/openid/web"
	"github.com/akatranlp/sentinel/utils"

	"github.com/go-viper/mapstructure/v2"
)

func (ip *IdentitiyProvider) OauthAuthorize(w http.ResponseWriter, r *http.Request) {
	var err error
	if r.URL.Fragment != "" {
		sendErrorPage(w, r, AuthorizeErrorTypeInvalidRequest.String(), "no fragment parameter allowed", http.StatusBadRequest)
		return
	}

	if err = r.ParseForm(); err != nil {
		sendErrorPage(w, r, AuthorizeErrorTypeInvalidRequest.String(), err.Error(), http.StatusBadRequest)
		return
	}

	clientID := r.FormValue(AuthorizeFormValueClientId.String())
	redirectURI, err := url.ParseRequestURI(r.FormValue(AuthorizeFormValueRedirectUri.String()))
	if err != nil {
		sendErrorPage(w, r, AuthorizeErrorTypeInvalidRequest.String(), err.Error(), http.StatusBadRequest)
		return
	}

	reg, ok := ip.clients[clientID]
	if !ok {
		sendErrorPage(w, r, AuthorizeErrorTypeInvalidRequest.String(), "invalid client id", http.StatusBadRequest)
		return
	}

	if err = reg.CheckRedirectURI(*redirectURI); err != nil {
		sendErrorPage(w, r, AuthorizeErrorTypeInvalidRequest.String(), "invalid redirect uri", http.StatusBadRequest)
		return
	}

	var formValues types.AuthRequest
	config := utils.CreateDecoderConfig()
	config.Result = &formValues
	decoder, err := mapstructure.NewDecoder(&config)
	if err != nil {
		handleClientError(w, r, formValues, AuthorizeError{
			AuthorizeErrorTypeInvalidRequest,
			err.Error(),
		})
		return
	}

	params, err := mapFormValueKeys(r.Form, maps.Keys(_AuthorizeFormValueValue))
	if err != nil {
		handleClientError(w, r, formValues, AuthorizeError{
			AuthorizeErrorTypeInvalidRequest,
			err.Error(),
		})
		return
	}

	if err = decoder.Decode(params); err != nil {
		handleClientError(w, r, formValues, AuthorizeError{
			AuthorizeErrorTypeInvalidRequest,
			err.Error(),
		})
		return
	}

	if err = formValues.IsValid(); err != nil {
		handleClientError(w, r, formValues, AuthorizeError{
			AuthorizeErrorTypeInvalidRequest,
			err.Error(),
		})
		return
	}

	ctx := r.Context()
	userID := ip.sessionManager.GetAuth(ctx)
	authTime := ip.sessionManager.GetAuthTime(ctx)
	isAuthed := userID != ""

	// If max age is set and our authtime is less than this
	if isAuthed && formValues.MaxAge != nil && time.Since(authTime) > (time.Duration(*formValues.MaxAge)*time.Second) {
		ip.sessionManager.SetAuth(ctx, "")
	} else if isAuthed && formValues.Prompt == enums.PromptLogin {
		ip.sessionManager.SetAuth(ctx, "")
	}

	// We need to ask the session again because the previous check could have logged out the user
	isAuthed = ip.sessionManager.IsAuthed(ctx)
	if !isAuthed && formValues.Prompt == enums.PromptNone {
		handleClientError(w, r, formValues, AuthorizeError{
			ErrorType: AuthorizeErrorTypeInteractionRequired,
		})
		return
	}

	if !isAuthed {
		ip.sessionManager.SetAuthRequest(ctx, formValues)
		http.Redirect(w, r, ip.basePath+"/login", http.StatusFound)
		return
	}

	ip.sendAuthResponse(w, r, AuthTokenValues{
		AuthRequest: formValues,
		UserID:      userID,
		AuthTime:    authTime,
	})
}

func (ip *IdentitiyProvider) sendAuthResponse(w http.ResponseWriter, r *http.Request, formValues AuthTokenValues, flashMessages ...templ.Component) {
	redirectURL := *formValues.RedirectURI
	params := make(url.Values)
	if formValues.State != "" {
		params.Set("state", formValues.State)
	}

	switch formValues.ResponseType {
	case enums.ResponseTypeCode:
		code := ip.createAuthToken(formValues)
		params.Set("code", code)
	}

	switch formValues.ResponseMode {
	case enums.ResponseModeFormPost:
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Pragma", "no-cache")
		web.FormPost(redirectURL.String(), params, flashMessages...).Render(r.Context(), w)
	case enums.ResponseModeFragment:
		redirectURL.Fragment = params.Encode()
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Pragma", "no-cache")
		web.FormRedirect(redirectURL.String(), flashMessages...).Render(r.Context(), w)
	case enums.ResponseModeQuery:
		redirectURL.RawQuery = params.Encode()
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Pragma", "no-cache")
		web.FormRedirect(redirectURL.String(), flashMessages...).Render(r.Context(), w)
	}
}

func sendErrorPage(w http.ResponseWriter, r *http.Request, errorType, error string, status int) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")
	w.WriteHeader(status)
	web.ErrorPage(errorType, error).Render(r.Context(), w)
}

func handleClientError(w http.ResponseWriter, r *http.Request, formValues types.AuthRequest, err AuthorizeError) {
	rURI := *formValues.RedirectURI

	params := make(url.Values)
	params.Set("error", err.ErrorType.String())
	if err.ErrorDescription != "" {
		params.Set("error_description", err.ErrorDescription)
	}
	if formValues.State != "" {
		params.Set("state", formValues.State)
	}

	switch formValues.ResponseMode {
	case enums.ResponseModeFormPost:
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Pragma", "no-cache")
		web.FormPost(rURI.String(), params).Render(r.Context(), w)
	case enums.ResponseModeFragment:
		rURI.Fragment = params.Encode()
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Pragma", "no-cache")
		web.FormRedirect(rURI.String()).Render(r.Context(), w)
	default:
		rURI.RawQuery = params.Encode()
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Pragma", "no-cache")
		web.FormRedirect(rURI.String()).Render(r.Context(), w)
	}
	http.Redirect(w, r, rURI.String(), http.StatusFound)
}

// ENUM(
// client_id,
// redirect_uri,
// scope,
// response_type,
// response_mode,
// state,
// nonce,
// code_challenge,
// code_challenge_method,
// )
type AuthorizeFormValue string

// ENUM(
//
//	// These are standard OAUTH 2.0 errorTypes
//	invalid_request, unauthorized_client, access_denied, unsupported_response_type, invalid_scope, server_error
//	// OpenID Connect ErrorTypes
//	interaction_required, login_required, account_selection_required, consent_required, invalid_request_uri, invalid_request_object, requets_not_supported, request_uri_not_supported, registration_not_supported
//
// )
type AuthorizeErrorType string

type AuthorizeError struct {
	ErrorType        AuthorizeErrorType `json:"error"`
	ErrorDescription string             `json:"error_description,omitzero"`
}
