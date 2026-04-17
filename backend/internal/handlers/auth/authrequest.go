package auth

import (
	"crypto/rand"
	"encoding/base64"
	"net/url"
	"stargazer/internal/core"
)

type AuthorizationRequest struct {
	State      string
	Nonce      string
	RequestUrl string
}

type AuthorizationRequestGenarator struct {
	config *core.AuthConfig
}

func NewAuthorizationRequestGenarator(config *core.AuthConfig) *AuthorizationRequestGenarator {
	return &AuthorizationRequestGenarator{
		config: config,
	}
}

func (g *AuthorizationRequestGenarator) Generate() *AuthorizationRequest {
	state := g.generateRandomString()
	nonce := g.generateRandomString()
	requestUrl := ""

	return &AuthorizationRequest{
		State:      state,
		Nonce:      nonce,
		RequestUrl: requestUrl,
	}
}

func (g *AuthorizationRequestGenarator) generateRandomString() string {
	bytes := make([]byte, 36)
	rand.Read(bytes)

	return base64.RawURLEncoding.EncodeToString(bytes)
}

func (g *AuthorizationRequestGenarator) generateRequestUrl(state string, nonce string) string {
	params := url.Values{}
	params.Set("client_id", g.config.ClientId)
	params.Set("redirect_uri", g.config.RedirectUri)
	params.Set("response_type", "code")
	params.Set("scope", "openid email")
	params.Set("state", state)
	params.Set("nonce", nonce)

	requestUrl := url.URL{
		Scheme:   "https",
		Host:     "accounts.google.com",
		Path:     "/o/oauth2/v2/auth",
		RawQuery: params.Encode(),
	}
	return requestUrl.String()
}
