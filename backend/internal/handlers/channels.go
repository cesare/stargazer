package handlers

import (
	"stargazer/internal/core"

	"github.com/gin-gonic/gin"
)

func RegisterChannelsHandler(group *gin.RouterGroup, appState *core.AppState) {
	group.POST("", func(c *gin.Context) {
	})
}
