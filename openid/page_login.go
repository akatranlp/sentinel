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
	"github.com/akatranlp/sentinel/openid/types"
	"github.com/akatranlp/sentinel/openid/web"
	"github.com/akatranlp/sentinel/provider"
	csrf "github.com/akatranlp/sentinel/session/gorilla_csrf"
	"github.com/akatranlp/sentinel/utils"
)

func (ip *IdentitiyProvider) LoginPage(w http.ResponseWriter, r *http.Request) {
	if ip.sessionManager.IsAuthed(r.Context()) {
		http.Redirect(w, r, ip.basePath+"/", http.StatusTemporaryRedirect)
		return
	}

	provs := slices.SortedFunc(its.Map21(maps.All(ip.providers), func(slug string, p provider.Provider) web.Provider {
		return web.Provider{
			Name:    p.GetName(),
			Slug:    slug,
			Icon:    string(p.GetType()),
			IconURL: utils.ParseIconURL(ip.basePath, p.GetIconURL()),
		}
	}), func(a, b web.Provider) int {
		return strings.Compare(a.Slug, b.Slug)
	})

	reactProvs := slices.SortedFunc(its.Map21(maps.All(ip.providers), func(slug string, p provider.Provider) types.Provider {
		return types.Provider{
			LoginURL:    fmt.Sprintf("%s/%s/login?redirect=%s", ip.basePath, slug, "/auth/"),
			Alias:       string(p.GetType()),
			ProviderID:  slug,
			DisplayName: p.GetName(),
			IconPath:    utils.ParseIconURL(ip.basePath, p.GetIconURL()),
		}
	}), func(a, b types.Provider) int {
		return strings.Compare(a.Alias, b.Alias)
	})

	sentinelCtx := types.NewSentinelCtx(ip.basePath, nil, nil)
	loginCtx := types.NewLoginSentinelCtx(sentinelCtx, reactProvs, types.CSRF{
		FieldName: ip.sessionManager.CsrfFormField(),
		Value:     csrf.Token(r),
	})

	w.Header().Set("Content-Type", "text/html")
	if err := ip.templates.ExecuteTemplate(w, "login.tmpl", loginCtx); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	web.Login(provs, ip.sessionManager.CsrfFormField(), csrf.Token(r), "/auth/").Render(r.Context(), io.Discard)
}
