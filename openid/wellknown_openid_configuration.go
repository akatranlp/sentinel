package openid

import (
	"encoding/json"
	"net/http"

	"github.com/akatranlp/sentinel/jose"
)

func (ip *IdentitiyProvider) WellKnownOpenIDConfiguration(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	j := jose.GetJose(ctx)
	openIDConfig := j.GetOpenIDConfiguration()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(openIDConfig)
}
