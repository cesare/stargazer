package models

import (
	"stargazer/internal/values"
	"time"
)

type Video struct {
	Id           values.VideoId   `db:"id"`
	ChannelId    values.ChannelId `db:"channel_id"`
	Title        string           `db:"title"`
	Description  string           `db:"description"`
	ThumbnailUrl string           `db:"thumbnail_url"`
	StartsAt     time.Time        `db:"starts_at"`
}
