package core

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
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

func (appState *AppState) AcquireDatabaseConnection(ctx context.Context) (*pgxpool.Conn, error) {
	conn, err := appState.dbpool.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to acquire database connection: %s", err)
	}

	return conn, nil
}

func (appState *AppState) BeginDatabaseTransaction(ctx context.Context) (pgx.Tx, error) {
	tx, err := appState.dbpool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin database transaction: %w", err)
	}

	return tx, nil
}
