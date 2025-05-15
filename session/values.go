package session

import (
	"context"
	"encoding/gob"
	"fmt"
	"time"

	"github.com/akatranlp/sentinel/openid/types"
)

const (
	csrfTokenKey = "csrfToken"
	userIDKey    = "userID"

	authKey     = "auth"
	authTimeKey = "auth-time"
	authReqKey  = "auth-req"

	providerVerifier = "verifier"
)

func (sm *SessionManager) getCsrfToken(ctx context.Context) string {
	return sm.GetString(ctx, csrfTokenKey)
}

func (sm *SessionManager) setCsrfToken(ctx context.Context, token string) {
	sm.Put(ctx, csrfTokenKey, token)
}

func (sm *SessionManager) SetAuth(ctx context.Context, id string) {
	if id == "" {
		sm.Remove(ctx, authKey)
		sm.Remove(ctx, authTimeKey)
	} else {
		sm.Put(ctx, authTimeKey, time.Now())
		sm.Put(ctx, authKey, id)
	}
	if err := sm.RenewToken(ctx); err != nil {
		fmt.Println(err)
	}
}

func (sm *SessionManager) GetAuth(ctx context.Context) string {
	return sm.GetString(ctx, authKey)
}

func (sm *SessionManager) GetAuthTime(ctx context.Context) time.Time {
	return sm.GetTime(ctx, authKey)
}

func (sm *SessionManager) IsAuthed(ctx context.Context) bool {
	return sm.GetAuth(ctx) != ""
}

func (sm *SessionManager) SetAuthRequest(ctx context.Context, authReq types.AuthRequest) {
	sm.Put(ctx, authReqKey, authReq)
}

func (sm *SessionManager) GetAuthRequest(ctx context.Context) (types.AuthRequest, bool) {
	authReq, ok := sm.Pop(ctx, authReqKey).(types.AuthRequest)
	return authReq, ok
}

func (sm *SessionManager) SetVerifier(ctx context.Context, verifier types.Verifier) {
	sm.Put(ctx, providerVerifier, verifier)
}

func (sm *SessionManager) GetVerifier(ctx context.Context) (types.Verifier, bool) {
	verifier, ok := sm.Pop(ctx, providerVerifier).(types.Verifier)
	return verifier, ok
}

func init() {
	gob.Register(types.AuthRequest{})
	gob.Register(types.Verifier{})
	gob.Register(time.Time{})
}
