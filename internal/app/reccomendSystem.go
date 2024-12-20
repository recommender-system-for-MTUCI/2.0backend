package app

import (
	"context"
	client2 "github.com/recommender-system-for-MTUCI/2.0backend/internal/client"
	"github.com/recommender-system-for-MTUCI/2.0backend/internal/config"
	"github.com/recommender-system-for-MTUCI/2.0backend/internal/controller"
	"github.com/recommender-system-for-MTUCI/2.0backend/internal/logger"
	storage "github.com/recommender-system-for-MTUCI/2.0backend/internal/storage"
	"github.com/recommender-system-for-MTUCI/2.0backend/internal/storage/postgres"
	"github.com/recommender-system-for-MTUCI/2.0backend/internal/token"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

func RecommendSystem() {
	ctx := context.Background()

	//init logger
	logCfg := logger.LoggerConfig()
	logger := zap.New(logCfg)
	logger.Info("Success to initialize logger")
	defer logger.Sync()

	//init config
	cfg, err := config.New()
	if err != nil {
		logger.Fatal("failed to initialize config", zap.Error(err))
		panic(err)
	}
	logger.Info("Success to initialize config")

	//init jwt
	jwt, err := token.NewToken(cfg)
	if err != nil {
		logger.Fatal("failed to initialize token", zap.Error(err))
	}
	logger.Info("Success to initialize jwt token")

	//init postgres
	pgx, err := postgres.New(ctx, logger, cfg)
	if err != nil {
		logger.Fatal("failed to initialize postgres", zap.Error(err))
		logger.Panic("failed to initialize postgres", zap.Error(err))
	}

	//init storage
	store, err := storage.NewStorage(logger, pgx)
	if err != nil {
		logger.Fatal("failed to initialize storage", zap.Error(err))
	}
	logger.Info("Success to initialize storage")

	//init client
	client, err := client2.New(logger, cfg)
	if err != nil {
		logger.Fatal("failed to initialize client", zap.Error(err))
	}
	logger.Info("Success to initialize client")

	//init server
	server := controller.New(logger, ctx, cfg, jwt, pgx, store, client)
	logger.Info("Success to initialize server")

	//run server
	go func() {
		logger.Info("Server run")
		err := server.Run()
		if err != nil {
			panic(err)
		}
	}()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	sign := <-stop
	logger.Info("Receive signal", zap.Any("signal", sign))

	//server shutdown
	server.Shutdown(ctx)
	logger.Info("Server shutdown")
}
