package session

import (
	// "context"
	// "fmt"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
)

type sessionManagerConfig struct {
	sessionName    string
	csrfFormField  string
	unAuthLifeTime time.Duration
	authLifeTime   time.Duration
	idleTimeout    time.Duration
}

var defaultConfig = sessionManagerConfig{
	sessionName:    "auth-session",
	csrfFormField:  "csrf-token",
	unAuthLifeTime: 7 * 24 * time.Hour,
	authLifeTime:   365 * 24 * time.Hour,
	idleTimeout:    7 * 24 * time.Hour,
}

type SessionManager struct {
	sessionManagerConfig
	*scs.SessionManager
}

func NewSessionManager(store scs.Store, basePath string, opts ...optionFn) *SessionManager {
	conf := defaultConfig
	for _, opt := range opts {
		opt(&conf)
	}

	sm := scs.New()
	sm.Store = store

	sm.Cookie.Domain = ""
	sm.Cookie.HttpOnly = true
	sm.Cookie.Secure = true
	sm.Cookie.Partitioned = true
	sm.Cookie.Path = basePath
	sm.Cookie.Name = conf.sessionName
	sm.Cookie.SameSite = http.SameSiteLaxMode
	sm.Lifetime = conf.authLifeTime
	sm.IdleTimeout = conf.idleTimeout

	return &SessionManager{
		sessionManagerConfig: conf,
		SessionManager:       sm,
	}
}
