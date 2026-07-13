package ytclient

import (
	"context"
	"fmt"
	"stargazer/internal/core"
	"stargazer/internal/values"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type VideoFinder struct {
	config *core.Config
}

func NewItemFinder(config *core.Config) *VideoFinder {
	return &VideoFinder{
		config: config,
	}
}

func (finder *VideoFinder) createService(ctx context.Context) (*youtube.Service, error) {
	option := option.WithAPIKey(finder.config.Youtube.ApiKey)
	return youtube.NewService(ctx, option)
}

func (finder *VideoFinder) createSearchService(ctx context.Context) (*youtube.SearchService, error) {
	service, err := finder.createService(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create youtube service: %s", err)
	}

	return youtube.NewSearchService(service), nil
}

func (finder *VideoFinder) createVideosService(ctx context.Context) (*youtube.VideosService, error) {
	service, err := finder.createService(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create youtube service: %s", err)
	}

	return youtube.NewVideosService(service), nil
}

func (finder *VideoFinder) ListUpcomingVideoIdsOnChannel(ctx context.Context, channelId values.ChannelId) ([]values.VideoId, error) {
	searchService, err := finder.createSearchService(ctx)
	if err != nil {
		return nil, err
	}

	call := searchService.List([]string{"snippet"}).EventType("upcoming").Type("video").ChannelId(fmt.Sprint(channelId))
	response, err := call.Do()
	if err != nil {
		return nil, fmt.Errorf("failed to list upcoming videos: %w", err)
	}

	var videoIds []values.VideoId
	for i := 0; i < len(response.Items); i++ {
		item := response.Items[i]
		videoIds = append(videoIds, values.VideoId(item.Id.VideoId))
	}

	return videoIds, nil
}

func (finder *VideoFinder) SearchVideos(ctx context.Context, videoIds []values.VideoId) ([]*youtube.Video, error) {
	ids := make([]string, len(videoIds))
	for i, v := range videoIds {
		ids[i] = fmt.Sprint(v)
	}

	videosService, err := finder.createVideosService(ctx)
	if err != nil {
		return nil, err
	}

	call := videosService.List([]string{"snippet", "liveStreamingDetails"}).Id(ids...)
	response, err := call.Do()
	if err != nil {
		return nil, fmt.Errorf("failed to search videos: %w", err)
	}

	return response.Items, nil
}
