package handlers

import (
	"net/http"
	"stargazer/internal/core"
	"stargazer/internal/handlers/auth"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func createOauth2Config(config *core.AuthConfig) *oauth2.Config {
	oauth2Config := oauth2.Config{
		ClientID:     config.ClientId,
		ClientSecret: config.ClientSecret,
		RedirectURL:  config.RedirectUri,
		Scopes: []string{
			"https://www.googleapis.com/auth/youtube.readonly",
		},
		Endpoint: google.Endpoint,
	}
	return &oauth2Config
}

func RegisterAuthHandlers(group *gin.RouterGroup, appState *core.AppState) {
	const SessionKeyAuthState = "google-auth-state"

	group.POST("", func(c *gin.Context) {
		state := auth.GenerateState()
		oauth2Config := createOauth2Config(&appState.Config.Auth)
		url := oauth2Config.AuthCodeURL(state)

		session := sessions.Default(c)
		session.Set(SessionKeyAuthState, state)
		session.Save()

		c.JSON(http.StatusOK, gin.H{
			"location": url,
		})
	})

	group.GET("/callback", func(c *gin.Context) {
		type callbackParams struct {
			Code  string `form:"code"`
			State string `form:"state"`
			Error string `form:"error"`
		}

		session := sessions.Default(c)
		savedState, stateOk := session.Get(SessionKeyAuthState).(string)
		session.Delete(SessionKeyAuthState)
		session.Save()

		if !stateOk {
			c.Status(http.StatusUnauthorized)
			return
		}

		var params callbackParams
		err := c.ShouldBindQuery(&params)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		if params.Error != "" {
			c.Status(http.StatusUnauthorized)
			return
		}

		if params.Code == "" || params.State == "" {
			c.Status(http.StatusUnauthorized)
			return
		}

		if params.State != savedState {
			c.Status(http.StatusUnauthorized)
			return
		}

		handleSuccessCallback(c, appState, params.Code)
	})
}

func handleSuccessCallback(c *gin.Context, appState *core.AppState, code string) {
	oauth2Config := createOauth2Config(&appState.Config.Auth)
	token, err := oauth2Config.Exchange(c, code)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
