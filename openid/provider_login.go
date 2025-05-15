package openid

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/akatranlp/sentinel/jose"
	"github.com/akatranlp/sentinel/openid/enums"
	"github.com/akatranlp/sentinel/openid/types"
	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

type oauthState struct {
	Redirect    string `json:"redirect"`
	RandomState string `json:"csrf"`
	State       string `json:"state"`
}

func (ip *identitiyProvider) ProviderLogin(w http.ResponseWriter, r *http.Request) {
	provider := r.PathValue("provider")

	p, ok := ip.providers[provider]
	if !ok {
		// TODO: Render nice frontend error-page
		http.Error(w, fmt.Sprintf("The provider %s is not configured", provider), http.StatusNotFound)
		return
	}

	clientState := r.FormValue("state")
	redirect := r.FormValue("redirect")

	verifier := types.Verifier{
		Verifier:    oauth2.GenerateVerifier(),
		RandomState: oauth2.GenerateVerifier(),
		Nonce:       oauth2.GenerateVerifier(),
	}

	state, err := json.Marshal(oauthState{State: clientState, RandomState: verifier.RandomState, Redirect: redirect})
	if err != nil {
		// TODO: Render nice frontend error-page
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx := r.Context()
	j := jose.GetJose(ctx)
	ip.sessionManager.SetVerifier(ctx, verifier)

	oauthConfig := p.GetOauthConfig(j.Issuer() + "/" + provider + "/callback")

	authCodeOption := oauth2.S256ChallengeOption(verifier.Verifier)
	responseMethod := oauth2.SetAuthURLParam("response_mode", enums.ResponseModeQuery.String())
	nonceOption := oidc.Nonce(verifier.Nonce)
	url := oauthConfig.AuthCodeURL(
		string(state),
		authCodeOption,
		responseMethod,
		nonceOption,
		oauth2.AccessTypeOffline,
	)

	http.Redirect(w, r, url, http.StatusFound)
}
