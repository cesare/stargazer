package core

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AppState struct {
	Config *Config
	dbpool *pgxpool.Pool
}

func CreateAppState(config *Config) (*AppState, error) {
	pool, err := pgxpool.New(context.Background(), config.Database.Url)
	if err != nil {
		return nil, fmt.Errorf("failed to establish database pool: %s", err.Error())
	}

	state := AppState{
		Config: config,
		dbpool: pool,
	}
	return &state, nil
}
