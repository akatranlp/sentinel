package openid

import (
	"maps"
	"net/http"
	"slices"
	"strings"

	"github.com/akatranlp/go-pkg/its"
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

	web.Login(provs, ip.sessionManager.CsrfFormField(), csrf.Token(r), "/auth/").Render(r.Context(), w)
}
