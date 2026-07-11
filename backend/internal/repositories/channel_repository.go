package repositories

import (
	"context"
	"fmt"
	"stargazer/internal/models"

	"github.com/jackc/pgx/v5"
)

type ChannelRepository struct {
	ctx        context.Context
	connection *pgx.Conn
}

func NewChannelRepository(ctx context.Context, conn *pgx.Conn) *ChannelRepository {
	return &ChannelRepository{
		ctx:        ctx,
		connection: conn,
	}
}

func (repository *ChannelRepository) Create(id string, title string, description string, thumbnailUrl string) (*models.Channel, error) {
	statement := "insert into channels (id, title, description, thumbnail_url) values ($1, $2, $3, $4) returning id, title, description, thumbnail_url"

	rows, err := repository.connection.Query(repository.ctx, statement, id, title, description, thumbnailUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to insert channel: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, fmt.Errorf("inserted channel missing")
	}

	var channel models.Channel
	err = rows.Scan(&channel.Id, &channel.Title, &channel.Description, &channel.ThumbnailUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to map row into channel struct: %w", err)
	}

	return &channel, nil
}

func (repository *ChannelRepository) TryFind(id string) (*models.Channel, error) {
	statement := "select id, title, description, thumbnail_url from channels where id = $1"
	rows, err := repository.connection.Query(repository.ctx, statement, id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch channel: %s, %w", id, err)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	var channel models.Channel
	err = rows.Scan(&channel.Id, &channel.Title, &channel.Description, &channel.ThumbnailUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to map result row into channel: %w", err)
	}

	return &channel, nil
}
