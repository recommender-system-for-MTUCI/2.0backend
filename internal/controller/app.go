package controller

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/recommender-system-for-MTUCI/2.0backend/internal/config"
	"go.uber.org/zap"
)

type Controller struct {
	logger *zap.Logger
	server *echo.Echo
	ctx    context.Context
	cfg    *config.Config
}

func New(logger *zap.Logger, ctx context.Context, cfg *config.Config) *Controller {
	return &Controller{
		logger: logger,
		server: echo.New(),
		ctx:    ctx,
		cfg:    cfg,
	}
}
func (ctrl *Controller) Run() error {
	ctrl.logger.Info("starting server")
	err := ctrl.server.Start(ctrl.cfg.Server.GetAddress())
	if err != nil {
		ctrl.logger.Error("failed to start server", zap.Error(err))
	}
	return nil
}
func (ctrl *Controller) Shutdown(ctx context.Context) error {
	ctrl.logger.Info("shutting down server")
	err := ctrl.server.Shutdown(ctx)
	if err != nil {
		ctrl.logger.Error("failed to shutdown server", zap.Error(err))
		panic(err)
	}
	return nil

}
