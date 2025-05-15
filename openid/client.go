package openid

import (
	"errors"
	"net/url"
	"slices"

	"github.com/akatranlp/sentinel/openid/enums"
)

type ClientRegistration struct {
	ClientID               string
	ClientSecret           string
	TokenExchangeSecret    string
	Scope                  enums.Scopes
	RedirectURIs           []string
	PostLogoutRedirectURIs []string
}

func (r *ClientRegistration) CheckRedirectURI(rURI url.URL) error {
	hostname := rURI.Hostname()
	switch hostname {
	case "localhost", "127.0.0.1", "::1":
		rURI.Host = hostname
	}

	if !slices.Contains(r.RedirectURIs, rURI.String()) {
		return errors.New("invalid redirect uri")
	}

	return nil
}

func (r *ClientRegistration) CheckPostLogoutRedirectURI(rURI url.URL) error {
	hostname := rURI.Hostname()
	switch hostname {
	case "localhost", "127.0.0.1", "::1":
		rURI.Host = hostname
	}

	if !slices.Contains(r.PostLogoutRedirectURIs, rURI.String()) {
		return errors.New("invalid redirect uri")
	}

	return nil
}
