package openid

import (
	"maps"
	"net/http"
	"slices"
	"strings"

	"github.com/akatranlp/go-pkg/its"
	"github.com/akatranlp/sentinel/account"
	"github.com/akatranlp/sentinel/openid/web"
	csrf "github.com/akatranlp/sentinel/session/gorilla_csrf"
)

func (ip *identitiyProvider) UserEditPage(w http.ResponseWriter, r *http.Request) {
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

	slices.SortFunc(accounts, func(a, b account.Account) int { return strings.Compare(a.Provider, b.Provider) })

	web.UserEdit(user, accounts, ip.sessionManager.CsrfFormField(), csrf.Token(r), "/auth/").Render(r.Context(), w)
}

func (ip *identitiyProvider) UserEdit(w http.ResponseWriter, r *http.Request) {
	if !ip.sessionManager.IsAuthed(r.Context()) {
		http.Redirect(w, r, ip.basePath+"/login", http.StatusTemporaryRedirect)
		return
	}

	if err := r.ParseForm(); err != nil {
		// TODO: send nice error page
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

	picture := r.FormValue("picture")
	name := r.FormValue("name")
	username := r.FormValue("username")
	email := r.FormValue("email")

	providers := maps.Collect(its.Map12(slices.Values(accounts), func(a account.Account) (string, account.Account) {
		return a.Provider, a
	}))

	newUser := user
	if p, ok := providers[picture]; ok {
		newUser.Picture = p.Picture
	} else {
		// TODO: send nice error page
		http.Error(w, "provider is not linked", http.StatusInternalServerError)
		return
	}

	if p, ok := providers[name]; ok {
		newUser.Name = p.Name
	} else {
		// TODO: send nice error page
		http.Error(w, "provider is not linked", http.StatusInternalServerError)
		return
	}
	if p, ok := providers[username]; ok {
		newUser.Username = p.PreferredUsername
	} else {
		// TODO: send nice error page
		http.Error(w, "provider is not linked", http.StatusInternalServerError)
		return
	}
	if p, ok := providers[email]; ok {
		newUser.Email = p.Email
	} else {
		// TODO: send nice error page
		http.Error(w, "provider is not linked", http.StatusInternalServerError)
		return
	}

	if err := ip.userStore.UpdateUser(r.Context(), user.UserID, newUser); err != nil {
		// TODO: send nice error page
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, ip.basePath+"/user", http.StatusFound)
}
