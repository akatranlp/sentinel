package store

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/akatranlp/sentinel/token"
)

type MemoryTokenStore struct {
	sessions map[string]token.Session
	mx       sync.Mutex
}

func (s *MemoryTokenStore) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"sessions": s.sessions,
	})
}

var _ (token.TokenStore) = (*MemoryTokenStore)(nil)

func NewMemoryTokenStore() *MemoryTokenStore {
	return &MemoryTokenStore{
		sessions: make(map[string]token.Session),
	}
}

func (s *MemoryTokenStore) SetSession(ctx context.Context, sid, jti string, expiry time.Time) error {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.sessions[sid] = token.Session{
		SessionID:  sid,
		RefreshJTI: jti,
		Expiry:     expiry,
	}
	return nil
}

func (s *MemoryTokenStore) GetSession(ctx context.Context, sid string) (token.Session, error) {
	s.mx.Lock()
	defer s.mx.Unlock()
	sess, ok := s.sessions[sid]
	if !ok {
		return token.Session{}, token.ErrSessionNotFound
	}
	if sess.Expiry.Before(time.Now()) {
		s.revokeSession(ctx, sid)
		return token.Session{}, token.ErrSessionNotFound
	}

	return sess, nil
}

func (s *MemoryTokenStore) revokeSession(ctx context.Context, sid string) error {
	delete(s.sessions, sid)
	return nil
}

func (s *MemoryTokenStore) RevokeSession(ctx context.Context, sid string) error {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.revokeSession(ctx, sid)
	return nil
}

func (s *MemoryTokenStore) StartSessionCleanup(ctx context.Context) {
	go func() {
		ticker := time.Tick(5 * time.Minute)
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker:
				curr := time.Now()
				s.mx.Lock()
				for sid, sess := range s.sessions {
					if sess.Expiry.Before(curr) {
						s.revokeSession(ctx, sid)
					}
				}
				s.mx.Unlock()
			}
		}
	}()
}
