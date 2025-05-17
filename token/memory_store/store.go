package store

import (
	"context"
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/akatranlp/sentinel/token"
)

type MemoryTokenStore struct {
	sessions map[string]token.Session
	mx       sync.Mutex
	savePath string
}

type marshal struct {
	Sessions map[string]token.Session `json:"sessions"`
}

func (s *MemoryTokenStore) MarshalJSON() ([]byte, error) {
	return json.Marshal(marshal{
		Sessions: s.sessions,
	})
}

func (s *MemoryTokenStore) UnmarshalJSON(data []byte) error {
	var store marshal

	if err := json.Unmarshal(data, &store); err != nil {
		return err
	}

	*s = MemoryTokenStore{
		sessions: store.Sessions,
	}
	return nil
}

var _ (token.TokenStore) = (*MemoryTokenStore)(nil)

func NewMemoryTokenStore(filePath ...string) (*MemoryTokenStore, error) {
	var path string
	if len(filePath) > 0 {
		path = filePath[0]
	}
	if path != "" {
		f, err := os.Open(path)
		if err == nil {
			defer f.Close()
			var store MemoryTokenStore
			if err := json.NewDecoder(f).Decode(&store); err != nil {
				return nil, err
			}
			store.savePath = path
			return &store, nil
		}
	}
	return &MemoryTokenStore{
		sessions: make(map[string]token.Session),
		savePath: path,
	}, nil
}

func (s *MemoryTokenStore) SetSession(ctx context.Context, sid, jti string, expiry time.Time) error {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.sessions[sid] = token.Session{
		SessionID:  sid,
		RefreshJTI: jti,
		Expiry:     expiry,
	}
	s.saveToFile()
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

func (s *MemoryTokenStore) revokeSession(_ context.Context, sid string) error {
	delete(s.sessions, sid)
	s.saveToFile()
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

func (s *MemoryTokenStore) saveToFile() error {
	if s.savePath == "" {
		return nil
	}
	f, err := os.Create(s.savePath)
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")
	return encoder.Encode(s)
}
