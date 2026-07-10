package handlers

import (
	"log/slog"
	"net/http"
	"net/url"
	"stargazer/internal/core"
	"stargazer/internal/repositories"
	"stargazer/internal/ytclient"
	"strings"

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

		url, err := url.ParseRequestURI(params.Url)
		if err != nil {
			slog.Warn("invalid url", "error", err, "url", url)
			c.Status(http.StatusBadRequest)
			return
		}

		pathElements := strings.Split(url.Path, "/")
		handle := pathElements[1]

		youtube := ytclient.NewChannelFinder(appState.Config)
		ch, err := youtube.FindByHandle(c, handle)
		if err != nil {
			slog.Error("failed in Youtube API FindByHandle", "error", err, "handle", handle)
			c.Status(http.StatusBadGateway)
			return
		}

		id := ch.Id
		title := ch.Snippet.Title
		description := ch.Snippet.Description
		thumbnailUrl := ch.Snippet.Thumbnails.Default.Url

		tx, err := appState.BeginDatabaseTransaction(c)
		if err != nil {
			slog.Error("failed to acquire database connection", "error", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		repository := repositories.NewChannelRepository(c, tx.Conn())
		channel, err := repository.TryFind(id)
		if err != nil {
			slog.Error("failed to find existing channel", "error", err, "channelId", id)
			tx.Rollback(c)
			c.Status(http.StatusInternalServerError)
			return
		}

		if channel == nil {
			channel, err = repository.Create(id, title, description, thumbnailUrl)
			if err != nil {
				slog.Error("failed to regisgter channel", "error", err, "channelId", id)
				tx.Rollback(c)
				c.Status(http.StatusInternalServerError)
				return
			}
		}
		tx.Commit(c)

		c.JSON(http.StatusCreated, gin.H{
			"id":           channel.Id,
			"title":        channel.Title,
			"description":  channel.Description,
			"thumbnailUrl": channel.ThumbnailUrl,
		})
	})
}
