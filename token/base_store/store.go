package basestore

import (
	"context"
	"errors"
	"time"

	"github.com/akatranlp/sentinel/token"
)

type Repository interface {
	GetSessionByID(ctx context.Context, sessionID string) (token.Session, error)
	UpdateSession(ctx context.Context, sess token.Session) error
	DeleteSessionByID(ctx context.Context, sessionID string) error
}

type SessionCreater interface {
	CreateSession(ctx context.Context, sess token.Session) (token.Session, error)
}

type SessionSetter interface {
	SetSession(ctx context.Context, sess token.Session) error
}

type SessionDeleter interface {
	DeleteSessionsAfterExpiry(ctx context.Context) error
}

type SessionGetter interface {
	GetAllSessions(ctx context.Context) ([]token.Session, error)
}

type BaseTokenStore struct {
	repo Repository
}

var _ (token.TokenStore) = (*BaseTokenStore)(nil)

func NewBaseTokenStore(repo Repository) (*BaseTokenStore, error) {
	_, ok1 := repo.(SessionDeleter)
	_, ok2 := repo.(SessionGetter)
	if !ok1 && !ok2 {
		return nil, errors.New("repo must be deleter or getter")
	}

	_, ok1 = repo.(SessionCreater)
	_, ok2 = repo.(SessionSetter)
	if !ok1 && !ok2 {
		return nil, errors.New("repo must be creater or setter")
	}

	return &BaseTokenStore{
		repo: repo,
	}, nil
}

func (s *BaseTokenStore) SetSession(ctx context.Context, sid, jti string, expiry time.Time) error {
	sess := token.Session{
		SessionID:  sid,
		RefreshJTI: jti,
		Expiry:     expiry,
	}
	if spec, ok := s.repo.(SessionSetter); ok {
		return spec.SetSession(ctx, sess)
	}

	spec := s.repo.(SessionCreater)

	_, err := s.repo.GetSessionByID(ctx, sid)
	if err == nil {
		return s.repo.UpdateSession(ctx, sess)
	} else if !errors.Is(err, token.ErrSessionNotFound) {
		return err
	}
	_, err = spec.CreateSession(ctx, sess)
	return err
}

func (s *BaseTokenStore) GetSession(ctx context.Context, sid string) (token.Session, error) {
	sess, err := s.repo.GetSessionByID(ctx, sid)
	if err != nil {
		return token.Session{}, err
	}
	if sess.Expiry.Before(time.Now()) {
		s.RevokeSession(ctx, sid)
		return token.Session{}, token.ErrSessionNotFound
	}

	return sess, nil
}

func (s *BaseTokenStore) RevokeSession(ctx context.Context, sid string) error {
	return s.repo.DeleteSessionByID(ctx, sid)
}

func (s *BaseTokenStore) StartSessionCleanup(ctx context.Context) error {
	go func() {
		ticker := time.Tick(5 * time.Minute)
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker:
				curr := time.Now()
				if spec, ok := s.repo.(SessionDeleter); ok {
					spec.DeleteSessionsAfterExpiry(ctx)
				} else if spec, ok := s.repo.(SessionGetter); ok {
					sessions, err := spec.GetAllSessions(ctx)
					if err != nil {
						continue
					}
					for _, sess := range sessions {
						if sess.Expiry.Before(curr) {
							s.RevokeSession(ctx, sess.SessionID)
						}
					}
				}
			}
		}
	}()
	return nil
}
