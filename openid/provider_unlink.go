package openid

import (
	"fmt"
	"net/http"

	"github.com/akatranlp/sentinel/account"
)

func (ip *IdentitiyProvider) ProviderUnlink(w http.ResponseWriter, r *http.Request) {
	provider := r.PathValue("provider")

	_, ok := ip.providers[provider]
	if !ok {
		// TODO: Render nice frontend error-page
		http.Error(w, fmt.Sprintf("The provider %s is not configured", provider), http.StatusNotFound)
		return
	}

	ctx := r.Context()

	if !ip.sessionManager.IsAuthed(ctx) {
		http.Redirect(w, r, ip.basePath+"/login", http.StatusTemporaryRedirect)
		return
	}

	userID := account.UserID(ip.sessionManager.GetAuth(ctx))
	account, err := ip.userStore.GetAccountByProvider(r.Context(), userID, provider)
	if err != nil {
		// TODO: send nice error page
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := ip.userStore.UnLinkAccount(r.Context(), userID, account.AccountID); err != nil {
		// TODO: send nice error page
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, ip.basePath+"/user", http.StatusFound)
}
