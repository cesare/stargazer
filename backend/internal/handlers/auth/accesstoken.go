package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"stargazer/internal/core"
)

type AccessTokenRequest struct {
	config *core.AuthConfig
}

type AccessTokenResponse struct {
	IdToken string `json:"id_token" binding:"required"`
}

func NewAccessTokenRequest(config *core.AuthConfig) *AccessTokenRequest {
	return &AccessTokenRequest{
		config: config,
	}
}

func (r *AccessTokenRequest) Execute(code string) (*AccessTokenResponse, error) {
	params := url.Values{}
	params.Set("client_id", r.config.ClientId)
	params.Set("client_secret", r.config.ClientSecret)
	params.Set("code", code)
	params.Set("grant_type", "authorization_code")
	params.Set("redirect_uri", r.config.RedirectUri)

	resp, err := http.PostForm("https://oauth2.googleapis.com/token", params)
	if err != nil {
		return nil, fmt.Errorf("access token request failed: %s", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read access token response: %s", err)
	}

	var response AccessTokenResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse access token response: %s", err)
	}

	return &response, nil
}
