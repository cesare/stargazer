package models

type Channel struct {
	Id           string `db:"id"`
	Title        string `db:"title"`
	Description  string `db:"description"`
	ThumbnailUrl string `db:"thumbnail_url"`
}
