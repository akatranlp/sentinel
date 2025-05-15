package session

import (
	"errors"
	"net/http"
	"strings"
	"time"

	csrf "github.com/akatranlp/sentinel/session/gorilla_csrf"
	"github.com/alexedwards/scs/v2"
)

type csrfStore struct {
	*SessionManager
}

var (
	ErrNoToken = errors.New("No Token in Session Storage")
)

func (s *csrfStore) Get(r *http.Request) ([]byte, error) {
	token := s.getCsrfToken(r.Context())
	if token == "" {
		return nil, ErrNoToken
	}
	return []byte(token), nil
}

func (s *csrfStore) Save(token []byte, r *http.Request, w http.ResponseWriter) error {
	s.SessionManager.setCsrfToken(r.Context(), strings.Clone(string(token)))
	return nil
}

func (sm *SessionManager) unAuthedSessionMaxTimeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if sm.IsAuthed(ctx) || sm.Status(r.Context()) == scs.Unmodified {
			next.ServeHTTP(w, r)
			return
		}
		expriy := sm.Deadline(r.Context())
		maxExpiry := time.Now().Add(sm.unAuthLifeTime).UTC()
		if expriy.After(maxExpiry) {
			sm.SetDeadline(r.Context(), maxExpiry)
		}
		next.ServeHTTP(w, r)
	})
}

func (sm *SessionManager) CsrfMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return csrf.Protect(
			&csrfStore{SessionManager: sm},
			csrf.FieldName(sm.csrfFormField),
		)(sm.unAuthedSessionMaxTimeMiddleware(next))
	}
}

func (sm *SessionManager) CsrfFormField() string {
	return sm.csrfFormField
}
