package server

import (
	"stargazer/internal/core"
	"stargazer/internal/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func createSessionMiddleware(config *core.Config) gin.HandlerFunc {
	store := cookie.NewStore(config.App.SessionKey.Bytes)
	return sessions.Sessions("stargazer-session", store)
}

func createCorsMiddleware(config *core.Config) gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins: []string{config.Frontend.BaseUrl},
		AllowMethods: []string{"DELETE", "GET", "OPTIONS", "POST"},
		AllowHeaders: []string{
			"Content-Type",
		},
		AllowCredentials: true,
	})
}

func CreateEngine(appState *core.AppState) *gin.Engine {
	engine := gin.Default()

	engine.Use(createSessionMiddleware(appState.Config))
	engine.Use(createCorsMiddleware(appState.Config))

	engine.GET("/ping", handlers.PingHandler)

	authGroup := engine.Group("/auth")
	handlers.RegisterAuthHandlers(authGroup, appState)

	return engine
}
