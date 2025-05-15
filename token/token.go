package token

import (
	"context"
	"errors"
	"time"
)

var (
	ErrSessionNotFound = errors.New("session not found")
)

type Session struct {
	SessionID  string
	RefreshJTI string
	Expiry     time.Time
}

type TokenStore interface {
	SetSession(ctx context.Context, sid string, jti string, expiry time.Time) error
	GetSession(ctx context.Context, sid string) (Session, error)
	RevokeSession(ctx context.Context, sid string) error
}
