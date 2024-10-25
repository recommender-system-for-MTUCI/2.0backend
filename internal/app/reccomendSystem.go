package app

import (
	"context"
	"github.com/recommender-system-for-MTUCI/2.0backend/internal/config"
	"github.com/recommender-system-for-MTUCI/2.0backend/internal/controller"
	"github.com/recommender-system-for-MTUCI/2.0backend/internal/logger"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

func ReccommendSystem() {
	ctx := context.Background()
	logCfg := logger.LoggerConfig()
	logger := zap.New(logCfg)
	defer logger.Sync()
	logger.Info("Logger initialize")
	cfg, err := config.New()
	if err != nil {
		logger.Error("Failed to initialize congif")
		panic(err)
	}
	logger.Info("Configuration initialize", zap.Any("config", cfg))
	server := controller.New(logger, ctx, cfg)
	logger.Info("Server initialize", zap.Any("server", server))
	go func() {
		logger.Info("Server run", zap.Any("server", server))
		err := server.Run()
		if err != nil {
			panic(err)
		}
	}()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	sign := <-stop
	logger.Info("Receive signal", zap.Any("signal", sign))
	server.Shutdown(ctx)
	logger.Info("Server shutdown")

}
