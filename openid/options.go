package openid

import (
	"io"
	"io/fs"
	"time"

	"github.com/akatranlp/sentinel/provider"
	"github.com/lestrrat-go/jwx/v3/jwk"
)

type OptionFn func(*ipConfig) error

// WithClients is the Option to create a clientRegistration.
// Use this multiple times to register more clients.
func WithClients(clients ...ClientRegistration) OptionFn {
	return func(ic *ipConfig) error {
		for _, client := range clients {
			ic.clients[client.ClientID] = client
		}
		return nil
	}
}

// WithProvider is the Option to register an Provider.
// Use this multiple times to register more provider.
func WithProviders(ps ...provider.Provider) OptionFn {
	return func(ic *ipConfig) error {
		for _, p := range ps {
			ic.providers[p.GetSlug()] = p
		}
		return nil
	}
}

func WithAccessTokenExpiration(expiration time.Duration) OptionFn {
	return func(ic *ipConfig) error {
		ic.atExpiration = expiration
		return nil
	}
}

func WithRefreshTokenExpiration(expiration time.Duration) OptionFn {
	return func(ic *ipConfig) error {
		ic.rtExpiration = expiration
		return nil
	}
}

func WithPublicURL(publicURL string) OptionFn {
	return func(ic *ipConfig) error {
		ic.publicURL = publicURL
		return nil
	}
}

func WithSessionUnAuthedLifeTime(lifeTime time.Duration) OptionFn {
	return func(ic *ipConfig) error {
		ic.sessionUnAuthedLifeTime = lifeTime
		return nil
	}
}

func WithSessionAuthedLifeTime(lifeTime time.Duration) OptionFn {
	return func(ic *ipConfig) error {
		ic.sessionAuthedLifeTime = lifeTime
		return nil
	}
}

func WithSessionIdleTimeout(timeout time.Duration) OptionFn {
	return func(ic *ipConfig) error {
		ic.sessionIdleTimout = timeout
		return nil
	}
}

func WithSessionName(name string) OptionFn {
	return func(ic *ipConfig) error {
		ic.sessionName = name
		return nil
	}
}

func WithSigningKey(sigingKey jwk.Key) OptionFn {
	return func(ic *ipConfig) error {
		ic.signingKey = sigingKey
		return nil
	}
}

func WithSigningKeyReader(r io.Reader) OptionFn {
	return func(ic *ipConfig) error {
		ic.signingKeyReader = r
		return nil
	}
}

func WithAppURL(appURL string) OptionFn {
	return func(ic *ipConfig) error {
		ic.appURL = appURL
		return nil
	}
}

func WithCustomAssetFS(fs fs.FS) OptionFn {
	return func(ic *ipConfig) error {
		ic.customAssets = fs
		return nil
	}
}
