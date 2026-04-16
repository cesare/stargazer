package server

import (
	"stargazer/internal/core"
	"stargazer/internal/handlers"

	"github.com/gin-gonic/gin"
)

func CreateEngine(appState *core.AppState) *gin.Engine {
	engine := gin.Default()

	engine.GET("/ping", handlers.PingHandler)

	return engine
}
