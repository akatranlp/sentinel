package memstore

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"sync"
	"time"
)

type item struct {
	object     []byte
	expiration int64
}

// MemStore represents the session store.
type MemStore struct {
	items       map[string]item
	mu          sync.RWMutex
	stopCleanup chan bool
	savePath    string
}

type marshal struct {
	Items map[string]marshalItem
}

type marshalItem struct {
	Object     string
	Expiration int64
}

func (s *MemStore) MarshalJSON() ([]byte, error) {
	items := make(map[string]marshalItem)
	for k, v := range s.items {
		items[k] = marshalItem{
			Object:     base64.StdEncoding.EncodeToString(v.object),
			Expiration: v.expiration,
		}
	}
	return json.Marshal(marshal{
		Items: items,
	})
}

func (s *MemStore) UnmarshalJSON(data []byte) error {
	var store marshal

	if err := json.Unmarshal(data, &store); err != nil {
		return err
	}

	items := make(map[string]item)
	for k, v := range store.Items {
		data, err := base64.StdEncoding.DecodeString(v.Object)
		if err != nil {
			return err
		}
		items[k] = item{
			object:     data,
			expiration: v.Expiration,
		}
	}

	*s = MemStore{
		items: items,
	}
	return nil
}

// New returns a new MemStore instance, with a background cleanup goroutine that
// runs every minute to remove expired session data.
func New(filePath ...string) (*MemStore, error) {
	return NewWithCleanupInterval(time.Minute)
}

// NewWithCleanupInterval returns a new MemStore instance. The cleanupInterval
// parameter controls how frequently expired session data is removed by the
// background cleanup goroutine. Setting it to 0 prevents the cleanup goroutine
// from running (i.e. expired sessions will not be removed).
func NewWithCleanupInterval(cleanupInterval time.Duration, filePath ...string) (*MemStore, error) {
	var path string
	if len(filePath) > 0 {
		path = filePath[0]
	}

	var s *MemStore

	if path != "" {
		f, err := os.Open(path)
		if err == nil {
			defer f.Close()
			var store MemStore
			if err := json.NewDecoder(f).Decode(&store); err != nil {
				return nil, err
			}
			store.savePath = path
			s = &store
		}
	}
	if s == nil {
		s = &MemStore{
			items:    make(map[string]item),
			savePath: path,
		}
	}

	if cleanupInterval > 0 {
		go s.startCleanup(cleanupInterval)
	}

	return s, nil
}

// Find returns the data for a given session token from the MemStore instance.
// If the session token is not found or is expired, the returned exists flag will
// be set to false.
func (m *MemStore) Find(token string) ([]byte, bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	item, found := m.items[token]
	if !found {
		return nil, false, nil
	}

	if time.Now().UnixNano() > item.expiration {
		return nil, false, nil
	}

	return item.object, true, nil
}

// Commit adds a session token and data to the MemStore instance with the given
// expiry time. If the session token already exists, then the data and expiry
// time are updated.
func (m *MemStore) Commit(token string, b []byte, expiry time.Time) error {
	m.mu.Lock()
	m.items[token] = item{
		object:     b,
		expiration: expiry.UnixNano(),
	}
	m.saveToFile()
	m.mu.Unlock()

	return nil
}

// Delete removes a session token and corresponding data from the MemStore
// instance.
func (m *MemStore) Delete(token string) error {
	m.mu.Lock()
	delete(m.items, token)
	m.saveToFile()
	m.mu.Unlock()

	return nil
}

// All returns a map containing the token and data for all active (i.e.
// not expired) sessions.
func (m *MemStore) All() (map[string][]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var mm = make(map[string][]byte)

	for token, item := range m.items {
		if item.expiration > time.Now().UnixNano() {
			mm[token] = item.object
		}
	}

	return mm, nil
}

func (m *MemStore) startCleanup(interval time.Duration) {
	m.stopCleanup = make(chan bool)
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ticker.C:
			m.deleteExpired()
		case <-m.stopCleanup:
			ticker.Stop()
			return
		}
	}
}

// StopCleanup terminates the background cleanup goroutine for the MemStore
// instance. It's rare to terminate this; generally MemStore instances and
// their cleanup goroutines are intended to be long-lived and run for the lifetime
// of your application.
//
// There may be occasions though when your use of the MemStore is transient.
// An example is creating a new MemStore instance in a test function. In this
// scenario, the cleanup goroutine (which will run forever) will prevent the
// MemStore object from being garbage collected even after the test function
// has finished. You can prevent this by manually calling StopCleanup.
func (m *MemStore) StopCleanup() {
	if m.stopCleanup != nil {
		m.stopCleanup <- true
	}
}

func (m *MemStore) deleteExpired() {
	now := time.Now().UnixNano()
	m.mu.Lock()
	for token, item := range m.items {
		if now > item.expiration {
			delete(m.items, token)
		}
	}
	m.saveToFile()
	m.mu.Unlock()
}

func (s *MemStore) saveToFile() error {
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
