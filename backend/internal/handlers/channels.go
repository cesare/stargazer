package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"stargazer/internal/core"
	. "stargazer/internal/errors"
	"stargazer/internal/handlers/channels"
	"stargazer/internal/views"

	"github.com/gin-gonic/gin"
)

func RegisterChannelsHandler(group *gin.RouterGroup, appState *core.AppState) {
	group.POST("", func(c *gin.Context) {
		type creationParams struct {
			Url string `json:"url" binding:"required"`
		}

		var params creationParams
		err := c.ShouldBind(&params)
		if err != nil {
			slog.Warn("invalid request", "error", err)
			c.Status(http.StatusBadRequest)
			return
		}

		registration := channels.NewRegistration(appState)
		channel, err := registration.Execute(c, params.Url)
		if err != nil {
			slog.Error("registration failed", "error", err, "requested-url", params.Url)

			var responseError ResponseError
			if errors.As(err, &responseError) {
				c.Status(responseError.Status())
			} else {
				c.Status(http.StatusInternalServerError)
			}
			return
		}

		view := views.NewChannelView(channel)
		c.JSON(http.StatusCreated, view)
	})
}
