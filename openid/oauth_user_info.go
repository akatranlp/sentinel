package openid

import (
	"encoding/json"
	"net/http"

	"github.com/akatranlp/sentinel/account"
	"github.com/akatranlp/sentinel/jose"
	"github.com/akatranlp/sentinel/openid/enums"
	"github.com/akatranlp/sentinel/openid/types"
	"github.com/akatranlp/sentinel/utils"
	"github.com/lestrrat-go/jwx/v3/jwt"
)

func (ip *identitiyProvider) OauthUserInfo(w http.ResponseWriter, r *http.Request) {
	j := jose.GetJose(r.Context())
	token, err := jwt.ParseHeader(
		r.Header, "Authorization",
		jwt.WithIssuer(j.Issuer()),
		jwt.WithKeySet(j.PublicKeys()),
		jwt.WithClaimValue(enums.ClaimTokenType.String(), enums.OauthTokenTypeAccessToken.String()),
	)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set("WWW-Authenticate", `Bearer realm="`+j.Issuer()+`, error="invalid_token", error_description="`+err.Error()+`"`)
		return
	}

	userID := utils.Bang(token.Subject())
	user, err := ip.userStore.GetUserByID(r.Context(), account.UserID(userID))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set("WWW-Authenticate", `Bearer realm="`+j.Issuer()+`, error="invalid_token", error_description="`+err.Error()+`"`)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(types.UserInfoResponse{
		Subject:           userID,
		Email:             user.Email,
		EmailVerified:     user.EmailVerified,
		Name:              user.Name,
		PreferredUsername: user.Username,
		Nickname:          user.Username,
		Picture:           user.Picture,
		Profile:           j.Issuer() + "/user",
	})

}
