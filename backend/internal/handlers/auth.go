package handlers

import (
	"net/http"
	"stargazer/internal/core"
	"stargazer/internal/handlers/auth"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func RegisterAuthHandlers(group *gin.RouterGroup, appState *core.AppState) {
	group.POST("", func(c *gin.Context) {
		generator := auth.NewAuthorizationRequestGenarator(&appState.Config.Auth)
		authRequest := generator.Generate()

		session := sessions.Default(c)
		session.Set("google-auth-state", authRequest.State)
		session.Set("google-auth-nonce", authRequest.Nonce)
		session.Save()

		c.JSON(http.StatusOK, gin.H{
			"location": authRequest.RequestUrl,
		})
	})

	group.POST("/callback", func(c *gin.Context) {
	})
}
