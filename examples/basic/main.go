package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"io"
	"log"
	"os"
	"time"

	"embed"
	usermemorystore "github.com/akatranlp/sentinel/account/memory_store"
	"github.com/akatranlp/sentinel/openid"
	"github.com/akatranlp/sentinel/openid/enums"
	tokenmemorystore "github.com/akatranlp/sentinel/token/memory_store"
	"github.com/akatranlp/sentinel/utils"
	"github.com/alexedwards/scs/v2/memstore"
	"github.com/lestrrat-go/jwx/v3/jwk"
)

//go:embed *.svg
var customAssets embed.FS

func GetOrCreateKey() (io.ReadCloser, error) {
	f, err := os.Open("examples/basic/jwk.json")
	if err == nil {
		return f, nil
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, err
	}
	jwkKey, err := jwk.Import(privateKey)
	if err != nil {
		return nil, err
	}
	if err = jwk.AssignKeyID(jwkKey); err != nil {
		return nil, err
	}
	if err = jwkKey.Set(jwk.KeyUsageKey, "sig"); err != nil {
		return nil, err
	}
	if err = jwkKey.Set(jwk.AlgorithmKey, "RS256"); err != nil {
		return nil, err
	}
	f, err = os.Create("examples/basic/jwk.json")
	if err != nil {
		return nil, err
	}
	if err = json.NewEncoder(f).Encode(jwkKey); err != nil {
		f.Close()
		return nil, err
	}
	f.Close()

	return os.Open("examples/basic/jwk.json")
}

func main() {
	ctx := context.Background()

	f, err := GetOrCreateKey()
	if err != nil {
		panic(err)
	}

	userStore, err := usermemorystore.NewMemoryUserStore("examples/basic/users.json")
	if err != nil {
		panic(err)
	}
	tokenStore, err := tokenmemorystore.NewMemoryTokenStore("examples/basic/tokens.json")
	if err != nil {
		panic(err)
	}

	tokenStore.StartSessionCleanup(ctx)

	sessionStore := memstore.NewWithCleanupInterval(1 * time.Second)

	ip, err := openid.NewIdentityProvider(
		"/",
		userStore,
		tokenStore,
		sessionStore,
		openid.WithClients(openid.ClientRegistration{
			ClientID:            "git-classrooms",
			ClientSecret:        "",
			TokenExchangeSecret: "supersecretsecret",
			Scope:               enums.ScopeValues(),
			PostLogoutRedirectURIs: []string{
				"http://localhost/callback",
				"http://localhost/",
				"https://oidcdebugger.com/debug",
			},
			RedirectURIs: []string{
				"http://localhost/callback",
				"http://localhost/",
				"https://oidcdebugger.com/debug",
			},
		}),
		openid.WithProviders(utils.Must(InitProviders())...),
		openid.WithSessionUnAuthedLifeTime(10*time.Hour),
		openid.WithSessionAuthedLifeTime(30*24*time.Hour),
		openid.WithSigningKeyReader(f),
		openid.WithCustomAssetFS(customAssets),
	)
	if err != nil {
		panic(err)
	}

	f.Close()

	log.Println("Listening on port 3000")
	ip.StartServer(ctx, ":3000")
}
