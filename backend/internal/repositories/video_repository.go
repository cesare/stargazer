package repositories

import (
	"context"
	"fmt"
	"stargazer/internal/models"
	"stargazer/internal/values"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/jackc/pgx/v5"
)

type VideoRepository struct {
	ctx        context.Context
	connection *pgx.Conn
}

func NewVideoRepository(ctx context.Context, conn *pgx.Conn) *VideoRepository {
	return &VideoRepository{
		ctx:        ctx,
		connection: conn,
	}
}

func (repository *VideoRepository) Create(newVideo *models.Video) (*models.Video, error) {
	statement := heredoc.Doc(`
		insert into videos (id, channel_id, title, description, thumbnail_url, starts_at)
		values ($1, $2, $3, $4, $5, $6)
		returning id, channel_id, title, description, thumbnail_url, starts_at
	`)

	rows, err := repository.connection.Query(
		repository.ctx,
		statement,
		newVideo.Id,
		newVideo.ChannelId,
		newVideo.Title,
		newVideo.Description,
		newVideo.ThumbnailUrl,
		newVideo.StartsAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to insert video: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, fmt.Errorf("inserted video missing")
	}

	var video models.Video
	err = rows.Scan(
		&video.Id,
		&video.ChannelId,
		&video.Title,
		&video.Description,
		&video.ThumbnailUrl,
		&video.StartsAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to map row into video struct: %w", err)
	}

	return &video, nil
}

func (repository *VideoRepository) Update(video *models.Video) (*models.Video, error) {
	statement := heredoc.Doc(`
		update videos
		set
			title = $2
			, description = $3
			, thumbnail_url = $4
			, starts_at = $5
		where
			id = $1
		returning id, channel_id, title, description, thumbnail_url, starts_at
	`)

	rows, err := repository.connection.Query(
		repository.ctx,
		statement,
		video.Id,
		video.Title,
		video.Description,
		video.ThumbnailUrl,
		video.StartsAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update video (id=%s): %w", video.Id, err)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, fmt.Errorf("updated video missing (id=%s)", video.Id)
	}

	var v models.Video
	err = rows.Scan(
		&v.Id,
		&v.ChannelId,
		&v.Title,
		&v.Description,
		&v.ThumbnailUrl,
		&v.StartsAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to map row into video: %w", err)
	}

	return &v, nil
}

func (repository *VideoRepository) Find(id values.VideoId) (*models.Video, error) {
	statement := heredoc.Doc(`
		select
			id
			, channel_id
			, title
			, description
			, thumbnail_url
			, starts_at
		from videos
		where
			id = $1
	`)

	rows, err := repository.connection.Query(repository.ctx, statement, id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch video (id=%s): %w", id, err)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	var v models.Video
	err = rows.Scan(
		&v.Id,
		&v.ChannelId,
		&v.Title,
		&v.Description,
		&v.ThumbnailUrl,
		&v.StartsAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to map row into video: %w", err)
	}

	return &v, nil
}

func (repository *VideoRepository) FindByIds(ids []values.VideoId) ([]models.Video, error) {
	statement := heredoc.Doc(`
		select
			id
			, channel_id
			, title
			, description
			, thumbnail_url
			, starts_at
		from videos
		where
			id = any($1)
	`)

	rows, err := repository.connection.Query(repository.ctx, statement, ids)
	if err != nil {
		return nil, fmt.Errorf("failed to lookup videos: %w", err)
	}
	defer rows.Close()

	var videos []models.Video
	for rows.Next() {
		var v models.Video
		err = rows.Scan(
			&v.Id,
			&v.ChannelId,
			&v.Title,
			&v.Description,
			&v.ThumbnailUrl,
			&v.StartsAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to map row into video: %w", err)
		}

		videos = append(videos, v)
	}

	return videos, nil
}
