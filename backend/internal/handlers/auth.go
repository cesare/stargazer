package handlers

import (
	"net/http"
	"stargazer/internal/core"
	"stargazer/internal/handlers/auth"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func RegisterAuthHandlers(group *gin.RouterGroup, appState *core.AppState) {
	const SessionKeyAuthState = "google-auth-state"
	const SessionKeyAuthNonce = "google-auth-nonce"

	group.POST("", func(c *gin.Context) {
		generator := auth.NewAuthorizationRequestGenarator(&appState.Config.Auth)
		authRequest := generator.Generate()

		session := sessions.Default(c)
		session.Set(SessionKeyAuthState, authRequest.State)
		session.Set(SessionKeyAuthNonce, authRequest.Nonce)
		session.Save()

		c.JSON(http.StatusOK, gin.H{
			"location": authRequest.RequestUrl,
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
		savedNonce, nonceOk := session.Get(SessionKeyAuthNonce).(string)
		session.Delete(SessionKeyAuthState)
		session.Delete(SessionKeyAuthNonce)
		session.Save()

		if !(stateOk && nonceOk) {
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

		handleSuccessCallback(appState, c, params.Code, savedNonce)
	})
}

func handleSuccessCallback(appState *core.AppState, c *gin.Context, code string, savedNonce string) {
}
