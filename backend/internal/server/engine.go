package server

import (
	"net/http"
	"stargazer/internal/core"

	"github.com/gin-gonic/gin"
)

func CreateEngine(appState *core.AppState) *gin.Engine {
	engine := gin.Default()

	engine.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	return engine
}
