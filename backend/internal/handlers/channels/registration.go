package channels

import (
	"context"
	"stargazer/internal/core"
	"stargazer/internal/errors"
	. "stargazer/internal/models"
	"stargazer/internal/repositories"
	"stargazer/internal/values"
	. "stargazer/internal/values"
	"stargazer/internal/ytclient"
)

type registration struct {
	appState *core.AppState
}

func NewRegistration(appState *core.AppState) *registration {
	return &registration{
		appState: appState,
	}
}

func (r *registration) Execute(ctx context.Context, url string) (*Channel, error) {
	handle, err := values.TryNewHandleFromUrl(url)
	if err != nil {
		return nil, errors.NewBadRequestError("failed to extract handle from URL %s: %w", url, err)
	}

	youtube := ytclient.NewChannelFinder(r.appState.Config)
	ch, err := youtube.FindByHandle(ctx, handle)
	if err != nil {
		return nil, errors.NewBadGatewayError("failed to fetch Youtube channel: %w", err)
	}

	id := ChannelId(ch.Id)
	title := ch.Snippet.Title
	description := ch.Snippet.Description
	thumbnailUrl := ch.Snippet.Thumbnails.Default.Url

	tx, err := r.appState.BeginDatabaseTransaction(ctx)
	if err != nil {
		return nil, errors.NewInternalServerError("failed to begin databsase transaction: %w", err)
	}

	repository := repositories.NewChannelRepository(ctx, tx.Conn())
	channel, err := repository.TryFind(id)
	if err != nil {
		tx.Rollback(ctx)
		return nil, errors.NewInternalServerError("failed to find channel (id=%s): %w", id, err)
	}

	if channel == nil {
		channel, err = repository.Create(id, title, description, thumbnailUrl)
		if err != nil {
			tx.Rollback(ctx)
			return nil, errors.NewInternalServerError("failed to insert channel (id=%s): %w", id, err)
		}
	}
	tx.Commit(ctx)

	return channel, nil
}
