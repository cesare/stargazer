package ytclient

import (
	"context"
	"fmt"
	"stargazer/internal/core"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type ChannelFinder struct {
	config *core.Config
}

func NewChannelFinder(config *core.Config) *ChannelFinder {
	return &ChannelFinder{
		config: config,
	}
}

func (finder *ChannelFinder) createService(ctx context.Context) (*youtube.ChannelsService, error) {
	option := option.WithAPIKey(finder.config.Youtube.ApiKey)
	service, err := youtube.NewService(ctx, option)
	if err != nil {
		return nil, fmt.Errorf("failed to create youtube service: %s", err)
	}

	return youtube.NewChannelsService(service), nil
}

func (finder *ChannelFinder) FindById(ctx context.Context, id string) (*youtube.Channel, error) {
	channelsService, err := finder.createService(ctx)
	if err != nil {
		return nil, err
	}

	call := channelsService.List([]string{"snippet", "contentDetails"}).Id(id)
	response, err := call.Do()
	if err != nil {
		return nil, fmt.Errorf("channels listing failed: %s", err)
	}

	return response.Items[0], nil
}

func (finder *ChannelFinder) FindByHandle(ctx context.Context, handle string) (*youtube.Channel, error) {
	channelsService, err := finder.createService(ctx)
	if err != nil {
		return nil, err
	}

	call := channelsService.List([]string{"snippet", "contentDetails"}).ForHandle(handle)
	response, err := call.Do()
	if err != nil {
		return nil, fmt.Errorf("channels listing failed: %s", err)
	}

	return response.Items[0], nil
}
