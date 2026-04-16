package server

import (
	"stargazer/internal/core"
	"stargazer/internal/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CreateEngine(appState *core.AppState) *gin.Engine {
	engine := gin.Default()

	engine.Use(cors.New(cors.Config{
		AllowOrigins: []string{appState.Config.Frontend.BaseUrl},
		AllowMethods: []string{"DELETE", "GET", "OPTIONS", "POST"},
		AllowHeaders: []string{
			"Content-Type",
		},
		AllowCredentials: true,
	}))

	engine.GET("/ping", handlers.PingHandler)

	return engine
}
