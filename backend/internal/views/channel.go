package views

import . "stargazer/internal/models"

type ChannelView struct {
	Id           ChannelId `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	ThumbnailUrl string    `json:"thumbnailUrl"`
}

func NewChannelView(channel *Channel) *ChannelView {
	return &ChannelView{
		Id:           channel.Id,
		Title:        channel.Title,
		Description:  channel.Description,
		ThumbnailUrl: channel.ThumbnailUrl,
	}
}
