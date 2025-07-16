package openid

import (
	"fmt"
	"io"
	"maps"
	"net/http"
	"os"
	"slices"
	"strings"

	"github.com/akatranlp/go-pkg/its"
	"github.com/akatranlp/sentinel/account"
	"github.com/akatranlp/sentinel/openid/types"
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

	reactProvs := slices.SortedFunc(its.Map21(maps.All(ip.providers), func(slug string, p provider.Provider) types.Provider {
		isLinked := slices.ContainsFunc(accounts, func(a account.Account) bool { return a.Provider == p.GetSlug() })
		action := "login"
		if isLinked {
			action = "unlink"
		}
		return types.Provider{
			LoginURL:    fmt.Sprintf("%s/%s/%s?redirect=%s", ip.basePath, slug, action, "/auth/"),
			Alias:       string(p.GetType()),
			ProviderID:  slug,
			DisplayName: p.GetName(),
			IconPath:    utils.ParseIconURL(ip.basePath, p.GetIconURL()),
			IsLinked:    isLinked,
		}
	}), func(a, b types.Provider) int {
		return strings.Compare(a.Alias, b.Alias)
	})

	reactAccounts := slices.Collect(its.Map(slices.Values(accounts), func(a account.Account) types.Account {
		return types.Account{
			Provider: a.Provider,
			Name:     a.Name,
			Email:    a.Email,
			Username: a.PreferredUsername,
			Picture:  a.Picture,
		}
	}))

	reactUser := types.User{
		ID:       string(user.UserID),
		Name:     user.Name,
		Username: user.Username,
		Picture:  user.Picture,
		Email:    user.Email,
	}

	sentinelCtx := types.NewSentinelCtx(ip.basePath, nil, nil, nil)
	userCtx := types.NewUserSentinelCtx(sentinelCtx, reactUser, reactAccounts, reactProvs, types.CSRF{
		FieldName: ip.sessionManager.CsrfFormField(),
		Value:     csrf.Token(r),
	})

	w.Header().Set("Content-Type", "text/html")
	if err := ip.templates.ExecuteTemplate(w, "user.tmpl", userCtx); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	web.User(user, provs, accounts, ip.sessionManager.CsrfFormField(), csrf.Token(r), "/auth/").Render(r.Context(), io.Discard)
}
