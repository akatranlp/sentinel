package csrf

import "net/http"

// store represents the session storage used for CSRF tokens.
type store interface {
	// Get returns the real CSRF token from the store.
	Get(*http.Request) ([]byte, error)
	// Save stores the real CSRF token in the store and writes a
	// cookie to the http.ResponseWriter.
	// For non-cookie stores, the cookie should contain a unique (256 bit) ID
	// or key that references the token in the backend store.
	// csrf.GenerateRandomBytes is a helper function for generating secure IDs.
	Save(token []byte, r *http.Request, w http.ResponseWriter) error
}
