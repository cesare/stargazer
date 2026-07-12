package models

type ChannelId string

type Channel struct {
	Id           ChannelId `db:"id"`
	Title        string    `db:"title"`
	Description  string    `db:"description"`
	ThumbnailUrl string    `db:"thumbnail_url"`
}
