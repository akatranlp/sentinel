package openid

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/akatranlp/sentinel/account"
	"github.com/akatranlp/sentinel/openid/types"
	"github.com/akatranlp/sentinel/openid/web"
	csrf "github.com/akatranlp/sentinel/session/gorilla_csrf"
)

func (ip *IdentitiyProvider) LogoutPage(w http.ResponseWriter, r *http.Request) {
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

	reactUser := types.User{
		ID:       string(user.UserID),
		Name:     user.Name,
		Username: user.Username,
		Picture:  user.Picture,
		Email:    user.Email,
	}

	sentinelCtx := types.NewSentinelCtx(ip.basePath, nil, nil, nil)
	userCtx := types.NewLogoutSentinelCtx(sentinelCtx, reactUser, "", "", types.CSRF{
		FieldName: ip.sessionManager.CsrfFormField(),
		Value:     csrf.Token(r),
	})

	w.Header().Set("Content-Type", "text/html")
	if err := ip.templates.ExecuteTemplate(w, "logout.tmpl", userCtx); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	web.Logout(ip.sessionManager.CsrfFormField(), csrf.Token(r), "", "").Render(r.Context(), io.Discard)
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
