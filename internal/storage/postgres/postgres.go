package postgres

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/recommender-system-for-MTUCI/2.0backend/internal/config"
	"go.uber.org/zap"
)

type Postgres struct {
	pgx    *pgxpool.Pool
	ctx    context.Context
	logger *zap.Logger
	cfg    *config.Config
}

func New(pgx *pgxpool.Pool, ctx context.Context, logger *zap.Logger, cfg *config.Config) (*Postgres, error) {
	postgres := &Postgres{
		pgx:    pgx,
		ctx:    ctx,
		logger: logger,
		cfg:    cfg,
	}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	config, err := pgxpool.ParseConfig(cfg.Postgres.GetAddressPostgres())
	if err != nil {
		logger.Fatal("failed to parse config for db", zap.Error(err))
	}
	pool, err := pgx.N
}
