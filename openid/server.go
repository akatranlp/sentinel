package openid

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

func (ip *identitiyProvider) StartServer(ctx context.Context, addr string) {
	r := chi.NewRouter()

	basePath := ip.basePath
	if basePath == "" {
		basePath = "/"
	}

	r.Mount(basePath, ip.Handler())

	if ip.basePath != "" {
		r.Handle("/", http.RedirectHandler(ip.basePath, http.StatusTemporaryRedirect))
	}

	server := http.Server{
		Handler: r,
		Addr:    addr,
	}

	context.AfterFunc(ctx, func() {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			// TODO: Logging package from context
			fmt.Println(err)
		}
	})

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		// TODO: Logging package from context
		fmt.Println(err)
	}
}
