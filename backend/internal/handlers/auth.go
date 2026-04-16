package handlers

import (
	"stargazer/internal/core"

	"github.com/gin-gonic/gin"
)

func RegisterAuthHandlers(group *gin.RouterGroup, appState *core.AppState) {
	group.POST("", func(c *gin.Context) {
	})

	group.POST("/callback", func(c *gin.Context) {
	})
}
