package openid

import (
	"fmt"
	"maps"
	"net/http"
	"net/url"

	"github.com/akatranlp/go-pkg/its"
	"github.com/akatranlp/go-pkg/middleware"
	"github.com/akatranlp/sentinel/jose"
	"github.com/akatranlp/sentinel/openid/web/assets"
	"github.com/akatranlp/sentinel/openid/web/shared"
	csrf "github.com/akatranlp/sentinel/session/gorilla_csrf"
	// chiMiddleware "github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func (ip *IdentitiyProvider) Handler() http.Handler {
	r := chi.NewRouter()

	var corsMiddleware *cors.Cors
	// corsMiddleware = cors.AllowAll()
	corsMiddleware = cors.New(cors.Options{
		// AllowedOrigins: []string{"*"},
		AllowOriginFunc: func(r *http.Request, origin string) bool {
			fmt.Println(origin)
			originURL, err := url.Parse(origin)
			if err != nil {
				return false
			}
			if origin == "https://jwt.io" {
				return true
			}
			return its.Any(maps.Values(ip.clients), func(client ClientRegistration) bool {
				return client.CheckRedirectURI(*originURL) == nil
			})
		},
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
		},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
	})

	r.Use(middleware.Forwarded(0))
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			proto := middleware.Proto(r)
			if proto == "http" {
				*r = *csrf.PlaintextHTTPRequest(r)
			}
			next.ServeHTTP(w, r)
		})
	})
	// r.Use(chiMiddleware.Logger)
	r.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			j, err := ip.joseBuilder.Build(middleware.Origin(r))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			ctx = jose.SetJose(ctx, j)
			ctx = shared.SetBasePath(ctx, ip.basePath)
			ctx = shared.SetAppURL(ctx, ip.appURL)
			ctx = shared.SetPath(ctx, r.URL.Path)

			*r = *r.WithContext(ctx)
			h.ServeHTTP(w, r)
		})
	})
	r.Use(ip.sessionManager.LoadAndSave)

	r.Group(func(r chi.Router) {
		r.Use(ip.sessionManager.CsrfMiddleware())

		r.Handle("/", http.RedirectHandler(ip.basePath+"/user", http.StatusTemporaryRedirect))
		r.Get("/login", ip.LoginPage)
		r.Post("/login", ip.LoginPage)
		r.Get("/logout", ip.LogoutPage)
		r.Post("/logout", ip.Logout)

		r.Get("/user", ip.UserPage)
		r.Get("/user/edit", ip.UserEditPage)
		r.Post("/user/edit", ip.UserEdit)
	})

	if ip.customAssets != nil {
		r.Handle("/dist/custom/*", http.StripPrefix(ip.basePath+"/dist/custom", http.FileServerFS(ip.customAssets)))
	}
	r.Handle("/dist/*", http.StripPrefix(ip.basePath+"/", http.FileServerFS(assets.Assets)))

	r.Route("/oauth", func(r chi.Router) {
		r.Use(corsMiddleware.Handler)

		r.Get("/authorize", ip.OauthAuthorize)
		r.Post("/authorize", ip.OauthAuthorize)

		r.Post("/token", ip.OauthToken)

		r.Post("/introspect", ip.OauthIntrospect)
		r.Post("/revoke", ip.OauthRevoke)
		r.Get("/userinfo", ip.OauthUserInfo)
		r.Post("/userinfo", ip.OauthUserInfo)
		//
		r.Group(func(r chi.Router) {
			r.Use(func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.FormValue(ip.sessionManager.CsrfFormField()) == "" {
						csrf.UnsafeSkipCheck(r)
					}
					next.ServeHTTP(w, r)
				})
			})
			r.Use(ip.sessionManager.CsrfMiddleware())
			r.Get("/logout", ip.OauthLogout)
			r.Post("/logout", ip.OauthLogout)
		})

		r.Get("/discovery/keys", ip.OauthDiscoveryKeys)
	})

	r.Route("/.well-known", func(r chi.Router) {
		r.Use(corsMiddleware.Handler)

		r.Get("/openid-configuration", ip.WellKnownOpenIDConfiguration)
		r.Get("/oauth-authorization-server", ip.WellKnownOpenIDConfiguration)
	})

	r.Route("/{provider}", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(ip.sessionManager.CsrfMiddleware())
			r.Post("/login", ip.ProviderLogin)
			r.Post("/link", nil) // Until now we use /login even for linking
			r.Post("/unlink", ip.ProviderUnlink)
		})
		r.Get("/callback", ip.ProviderCallback)
		r.Post("/callback", ip.ProviderCallback)
	})

	return r
}
