package openid

import (
	"io"
	"io/fs"
	"sync"
	"time"

	"github.com/akatranlp/sentinel/account"
	"github.com/akatranlp/sentinel/jose"
	"github.com/akatranlp/sentinel/provider"
	"github.com/akatranlp/sentinel/session"
	"github.com/akatranlp/sentinel/token"
	"github.com/alexedwards/scs/v2"
	"github.com/lestrrat-go/jwx/v3/jwk"
)

type ipConfig struct {
	// Clients
	clients map[string]ClientRegistration

	// Provider
	providers map[string]provider.Provider

	// SessionManager
	sessionName             string
	sessionUnAuthedLifeTime time.Duration
	sessionAuthedLifeTime   time.Duration
	sessionIdleTimout       time.Duration

	// JOSE
	basePath         string
	publicURL        string
	atExpiration     time.Duration
	rtExpiration     time.Duration
	signingKey       jwk.Key
	signingKeyReader io.Reader

	// Web
	appURL       string
	customAssets fs.FS
}

var defaultConfig = ipConfig{
	clients:   make(map[string]ClientRegistration),
	providers: make(map[string]provider.Provider),
	basePath:  "/",
	publicURL: "",

	sessionName:             "auth-session",
	sessionUnAuthedLifeTime: 30 * time.Minute,
	sessionAuthedLifeTime:   364 * 24 * time.Hour,
	sessionIdleTimout:       7 * 24 * time.Hour,
}

type IdentitiyProvider struct {
	ipConfig
	joseBuilder    *jose.JoseBuilder
	userStore      account.UserStore
	tokenStore     token.TokenStore
	sessionManager *session.SessionManager
	authMap        sync.Map
}

func NewIdentityProvider(
	basePath string,
	userStore account.UserStore,
	tokenStore token.TokenStore,
	sessionStore scs.Store,
	opts ...OptionFn,
) (*IdentitiyProvider, error) {
	var err error
	conf := defaultConfig
	for _, opt := range opts {
		err = opt(&conf)
		if err != nil {
			return nil, err
		}
	}

	joseOpts := []jose.OptionFn{
		jose.WithBasePath(basePath),
		jose.WithPublicURL(conf.publicURL),
		jose.WithAccessTokenExpiration(conf.atExpiration),
		jose.WithRefeshTokenExpiration(conf.rtExpiration),
	}
	if conf.signingKey != nil {
		joseOpts = append(joseOpts, jose.WithSigningKey(conf.signingKey))
	}
	if conf.signingKeyReader != nil {
		joseOpts = append(joseOpts, jose.WithSigningKeyReader(conf.signingKeyReader))
	}
	jb, err := jose.NewJoseBuilder(joseOpts...)
	if err != nil {
		return nil, err
	}

	j, err := jb.Build("")
	if err != nil {
		return nil, err
	}
	conf.basePath = j.BasePath()

	sm := session.NewSessionManager(
		sessionStore,
		conf.basePath,
		session.WithSessionName(conf.sessionName),
		session.WithAuthLifeTime(conf.sessionAuthedLifeTime),
		session.WithUnAuthLifeTime(conf.sessionUnAuthedLifeTime),
		session.WithIdleTimeout(conf.sessionIdleTimout),
	)

	return &IdentitiyProvider{
		ipConfig:       conf,
		joseBuilder:    jb,
		userStore:      userStore,
		tokenStore:     tokenStore,
		sessionManager: sm,
		authMap:        sync.Map{},
	}, nil
}

func (ip *IdentitiyProvider) PublicKeys() jwk.Set {
	jose, _ := ip.joseBuilder.Build("")
	return jose.PublicKeys()
}
