package postgres

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/recommender-system-for-MTUCI/2.0backend/internal/config"
	"go.uber.org/zap"
)

func New(ctx context.Context, logger *zap.Logger, cfg *config.Config) (*pgxpool.Pool, error) {
	var pgx *pgxpool.Pool
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	addres, err := pgxpool.ParseConfig(cfg.Postgres.GetAddressPostgres())
	if err != nil {
		logger.Fatal("failed to parse config for db", zap.Error(err))
	}
	pgx, err = pgxpool.NewWithConfig(ctx, addres)
	if err != nil {
		logger.Fatal("failed to connect to db", zap.Error(err))
	}
	logger.Info("connected to db")
	return pgx, nil
}
