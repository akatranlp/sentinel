package openid

import (
	"encoding/json"
	"net/http"

	"github.com/akatranlp/sentinel/jose"
)

func (ip *IdentitiyProvider) OauthDiscoveryKeys(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	j := jose.GetJose(ctx)
	keys := j.PublicKeys()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(keys)
}
