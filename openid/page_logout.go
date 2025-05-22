package openid

import (
	"net/http"

	"github.com/akatranlp/sentinel/openid/web"
	csrf "github.com/akatranlp/sentinel/session/gorilla_csrf"
)

func (ip *IdentitiyProvider) LogoutPage(w http.ResponseWriter, r *http.Request) {
	if !ip.sessionManager.IsAuthed(r.Context()) {
		http.Redirect(w, r, ip.basePath+"/login", http.StatusTemporaryRedirect)
		return
	}

	web.Logout(ip.sessionManager.CsrfFormField(), csrf.Token(r), "", "").Render(r.Context(), w)
}

func (ip *IdentitiyProvider) Logout(w http.ResponseWriter, r *http.Request) {
	if !ip.sessionManager.IsAuthed(r.Context()) {
		http.Redirect(w, r, ip.basePath+"/login", http.StatusFound)
		return
	}

	if err := ip.sessionManager.Destroy(r.Context()); err != nil {
		// TODO: Render errorPage
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, ip.basePath+"/login", http.StatusFound)
}
