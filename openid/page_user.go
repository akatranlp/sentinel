package openid

import (
	"maps"
	"net/http"
	"slices"
	"strings"

	"github.com/akatranlp/go-pkg/its"
	"github.com/akatranlp/sentinel/account"
	"github.com/akatranlp/sentinel/openid/web"
	"github.com/akatranlp/sentinel/provider"
	csrf "github.com/akatranlp/sentinel/session/gorilla_csrf"
	"github.com/akatranlp/sentinel/utils"
)

func (ip *IdentitiyProvider) UserPage(w http.ResponseWriter, r *http.Request) {
	if !ip.sessionManager.IsAuthed(r.Context()) {
		http.Redirect(w, r, ip.basePath+"/login", http.StatusTemporaryRedirect)
		return
	}

	userID := account.UserID(ip.sessionManager.GetAuth(r.Context()))
	user, err := ip.userStore.GetUserByID(r.Context(), userID)
	if err != nil {
		// TODO: send nice error page
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	accounts, err := ip.userStore.GetAccountsForUserID(r.Context(), userID)
	if err != nil {
		// TODO: send nice error page
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	provs := slices.SortedFunc(its.Map21(maps.All(ip.providers), func(slug string, p provider.Provider) web.LinkProvider {
		return web.LinkProvider{
			Name:    p.GetName(),
			Slug:    slug,
			Icon:    string(p.GetType()),
			IconURL: utils.ParseIconURL(ip.basePath, p.GetIconURL()),
			Linked: slices.ContainsFunc(accounts, func(acc account.Account) bool {
				return acc.Provider == slug
			}),
		}
	}), func(a, b web.LinkProvider) int {
		return strings.Compare(a.Slug, b.Slug)
	})

	slices.SortFunc(accounts, func(a, b account.Account) int { return strings.Compare(a.Provider, b.Provider) })

	web.User(user, provs, accounts, ip.sessionManager.CsrfFormField(), csrf.Token(r), "/auth/").Render(r.Context(), w)
}
