package jose

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/akatranlp/sentinel/openid/enums"
	"github.com/akatranlp/sentinel/openid/types"

	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/lestrrat-go/jwx/v3/jwt"
	"github.com/lestrrat-go/jwx/v3/jwt/openid"
)

var (
	ErrInvalidOrigin   = errors.New("invalid origin provided")
	ErrInvalidBasePath = errors.New("invalid basePath provided")
)

type joseConfig struct {
	origin          string
	basePath        string
	accessTokenExp  time.Duration
	refreshTokenExp time.Duration

	signingKey jwk.Key
	publicKey  jwk.Key
}

type OptionFn func(*joseConfig) error

func WithBasePath(basePath string) OptionFn {
	return func(c *joseConfig) error {
		baseURL, err := url.Parse(basePath)
		if err != nil {
			return err
		}
		basePath = baseURL.EscapedPath()
		if !path.IsAbs(basePath) {
			return ErrInvalidBasePath
		} else if len(basePath) > 0 && basePath[len(basePath)-1] == '/' {
			basePath = basePath[:len(basePath)-1]
		}
		c.basePath = basePath
		return nil
	}
}

func WithPublicURL(origin string) OptionFn {
	return func(c *joseConfig) (err error) {
		c.origin, err = parseOrigin(origin)
		return err
	}
}

func WithAccessTokenExpiration(duration time.Duration) OptionFn {
	return func(c *joseConfig) error {
		c.accessTokenExp = duration
		return nil
	}
}

func WithRefeshTokenExpiration(duration time.Duration) OptionFn {
	return func(c *joseConfig) error {
		c.refreshTokenExp = duration
		return nil
	}
}

func WithSigningKey(signingKey jwk.Key) OptionFn {
	return func(jc *joseConfig) error {
		jc.signingKey = signingKey
		return nil
	}
}

func WithSigningKeyReader(signingKeyReader io.Reader) OptionFn {
	return func(jc *joseConfig) error {
		data, err := io.ReadAll(signingKeyReader)
		if err != nil {
			return err
		}
		signingKey, err := jwk.ParseKey(data)
		if err != nil {
			return err
		}
		jc.signingKey = signingKey
		return nil
	}
}

type jose struct {
	joseConfig
}

type JoseBuilder struct {
	joseConfig
}

func parseOrigin(origin string) (string, error) {
	if origin == "" {
		return "", nil
	}
	originURL, err := url.Parse(origin)
	if err != nil {
		return "", err
	}
	if originURL.RawFragment != "" {
		return "", ErrInvalidOrigin
	}
	if originURL.RawQuery != "" {
		return "", ErrInvalidOrigin
	}
	if originURL.RawPath != "" {
		return "", ErrInvalidOrigin
	}
	return strings.Clone(origin), nil
}

func (b *JoseBuilder) Build(origin string) (*jose, error) {
	var err error
	conf := b.joseConfig
	if origin != "" {
		if conf.origin, err = parseOrigin(origin); err != nil {
			return nil, err
		}
	}
	return &jose{
		joseConfig: conf,
	}, nil
}

var defaultConf = joseConfig{
	basePath:        "",
	accessTokenExp:  15 * time.Minute,
	refreshTokenExp: 7 * 24 * time.Hour,
}

func NewJoseBuilder(opts ...OptionFn) (*JoseBuilder, error) {
	var err error
	conf := defaultConf
	for _, opt := range opts {
		if err = opt(&conf); err != nil {
			return nil, err
		}
	}

	if conf.signingKey == nil {
		return nil, errors.New("no sigingKey specified")
	}

	conf.publicKey, err = conf.signingKey.PublicKey()
	if err != nil {
		return nil, err
	}
	return &JoseBuilder{
		joseConfig: conf,
	}, nil
}

func (j *jose) BasePath() string {
	return j.basePath
}

func (j *jose) Issuer() string {
	return j.origin + j.basePath
}

func (j *jose) createBaseToken(
	sub string,
	aud []string,
	sid string,
	scope enums.Scopes,
	currTime time.Time,
	expiration time.Duration,
) *jwt.Builder {
	return jwt.NewBuilder().
		Subject(sub).
		Issuer(j.Issuer()).
		IssuedAt(currTime).
		Audience(aud).
		Claim(enums.ClaimSid.String(), sid).
		Expiration(currTime.Add(expiration)).
		NotBefore(currTime).
		JwtID(uuid.NewString()).
		Claim(enums.ClaimScope.String(), scope.String())
}

func (j *jose) signToken(token jwt.Token, key jwk.Key) (string, error) {
	tokenBytes, err := jwt.Sign(token, jwt.WithKey(jwa.RS256(), key))
	if err != nil {
		return "", err
	}

	return string(tokenBytes), nil
}

type TokenCreateArg struct {
	CurrTime  time.Time
	Subject   string
	Audience  []string
	SessionID string
	Scope     enums.Scopes
}

func (j *jose) CreateAccessToken(arg TokenCreateArg) (token jwt.Token, signedToken string, err error) {
	token, err = j.createBaseToken(
		arg.Subject,
		arg.Audience,
		arg.SessionID,
		arg.Scope,
		arg.CurrTime,
		j.accessTokenExp,
	).
		Claim(enums.ClaimTokenType.String(), enums.OauthTokenTypeAccessToken.String()).
		Build()
	if err != nil {
		return nil, "", err
	}

	signedToken, err = j.signToken(token, j.signingKey)
	if err != nil {
		return nil, "", err
	}

	return token, signedToken, nil
}

func (j *jose) CreateRefreshToken(arg TokenCreateArg) (token jwt.Token, signedToken string, err error) {
	token, err = j.createBaseToken(
		arg.Subject,
		arg.Audience,
		arg.SessionID,
		arg.Scope,
		arg.CurrTime,
		j.refreshTokenExp,
	).
		Claim(enums.ClaimTokenType.String(), enums.OauthTokenTypeRefreshToken.String()).
		Build()
	if err != nil {
		return nil, "", err
	}

	signedToken, err = j.signToken(token, j.signingKey)
	if err != nil {
		return nil, "", err
	}

	return token, signedToken, nil
}

type IDTokenCreateArg struct {
	TokenCreateArg
	SignedAccessToken string
	Nonce             string
	AuthTime          time.Time
	Email             string
	EmailVerified     bool
	Name              string
	Username          string
	Picture           string
}

func (j *jose) CreateIDToken(arg IDTokenCreateArg) (token openid.Token, signedToken string, err error) {
	atSha := sha256.Sum256([]byte(arg.SignedAccessToken))
	atHash := base64.RawURLEncoding.EncodeToString(atSha[:16])

	token, err = openid.NewBuilder().
		Subject(arg.Subject).
		Issuer(j.Issuer()).
		IssuedAt(arg.CurrTime).
		Audience(arg.Audience).
		Claim(enums.ClaimSid.String(), arg.SessionID).
		Expiration(arg.CurrTime.Add(j.accessTokenExp)).
		NotBefore(arg.CurrTime).
		JwtID(uuid.NewString()).
		Claim(enums.ClaimScope.String(), arg.Scope.String()).
		Claim(enums.ClaimAtHash.String(), atHash).
		Claim(enums.ClaimNonce.String(), arg.Nonce).
		Claim(enums.ClaimAuthTime.String(), arg.AuthTime.Unix()).
		Email(arg.Email).
		EmailVerified(arg.EmailVerified).
		Name(arg.Name).
		Nickname(arg.Username).
		PreferredUsername(arg.Username).
		Profile(j.Issuer()+"/user").
		Picture(arg.Picture).
		Claim(enums.ClaimTokenType.String(), enums.OauthTokenTypeIdToken.String()).
		Build()
	if err != nil {
		return nil, "", err
	}

	signedToken, err = j.signToken(token, j.signingKey)
	if err != nil {
		return nil, "", err
	}

	return token, signedToken, nil
}

type Tokens struct {
	AccessToken        jwt.Token
	RefreshToken       jwt.Token
	IDToken            openid.Token
	SignedAccessToken  string
	SignedRefreshToken string
	SignedIDToken      string
}

func (j *jose) CreateTokens(arg IDTokenCreateArg, createRefreshToken, createIDToken bool) (*Tokens, error) {
	var tokens Tokens
	var err error
	tokens.AccessToken, tokens.SignedAccessToken, err = j.CreateAccessToken(arg.TokenCreateArg)
	if err != nil {
		return nil, err
	}
	if createRefreshToken {
		tokens.RefreshToken, tokens.SignedRefreshToken, err = j.CreateRefreshToken(arg.TokenCreateArg)
		if err != nil {
			return nil, err
		}
	}
	if createIDToken {
		arg.SignedAccessToken = tokens.SignedAccessToken
		tokens.IDToken, tokens.SignedIDToken, err = j.CreateIDToken(arg)
		if err != nil {
			return nil, err
		}
	}
	return &tokens, nil
}

func (j *jose) PublicKeys() jwk.Set {
	set := jwk.NewSet()
	_ = set.AddKey(j.publicKey)
	return set
}

func (j *jose) GetOpenIDConfiguration() types.OpenIDConfiguration {
	return types.CreateOIDCConfig(j.Issuer())
}
