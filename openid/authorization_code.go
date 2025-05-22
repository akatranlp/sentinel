package openid

import (
	"time"

	"github.com/akatranlp/sentinel/openid/types"
	"github.com/google/uuid"
)

type AuthTokenValues struct {
	types.AuthRequest
	UserID   string
	AuthTime time.Time
}

func (ip *IdentitiyProvider) createAuthToken(authReq AuthTokenValues) string {
	code := uuid.NewString()
	ip.authMap.Store(code, authReq)
	time.AfterFunc(10*time.Second, func() {
		ip.authMap.Delete(code)
	})
	return code
}

func (ip *IdentitiyProvider) validateCode(code string) (AuthTokenValues, bool) {
	v, ok := ip.authMap.LoadAndDelete(code)
	if !ok {
		return AuthTokenValues{}, false
	}
	return v.(AuthTokenValues), true
}
