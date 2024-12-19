package controller

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/recommender-system-for-MTUCI/2.0backend/internal/client"
	"github.com/recommender-system-for-MTUCI/2.0backend/internal/config"
	"github.com/recommender-system-for-MTUCI/2.0backend/internal/storage"
	token "github.com/recommender-system-for-MTUCI/2.0backend/internal/token"
	"go.uber.org/zap"
)

type Controller struct {
	logger  *zap.Logger
	server  *echo.Echo
	ctx     context.Context
	cfg     *config.Config
	token   token.JWT
	pgx     *pgxpool.Pool
	storage *storage.Storage
	client  *client.Client
}

// func to create new instance controller
func New(logger *zap.Logger, ctx context.Context, cfg *config.Config, token token.JWT, pgx *pgxpool.Pool, store *storage.Storage, client *client.Client) *Controller {

	ctrl := &Controller{
		logger:  logger,
		server:  echo.New(),
		ctx:     ctx,
		cfg:     cfg,
		token:   token,
		pgx:     pgx,
		storage: store,
		client:  client,
	}
	ctrl.RegisterMiddlewares()
	ctrl.RegisterRoutes()
	return ctrl
}

// func to run server
func (ctrl *Controller) Run() error {
	err := ctrl.server.Start("localhost:8080")
	if err != nil {
		ctrl.logger.Error("failed to start server", zap.Error(err))
	}
	return nil
}

// func to shutdown server
func (ctrl *Controller) Shutdown(ctx context.Context) error {
	ctrl.logger.Info("shutting down server")
	err := ctrl.server.Shutdown(ctx)
	if err != nil {
		ctrl.logger.Error("failed to shutdown server", zap.Error(err))
		panic(err)
	}
	return nil

}

// func to configure routes with handlers
func (ctrl *Controller) RegisterRoutes() {
	ctrl.logger.Info("registering routes")
	api := ctrl.server.Group("/api")
	api.GET("/recommend_system", ctrl.handleGetMainPage)
	api.POST("/registration", ctrl.handleRegistration)
	api.GET("/favorites", ctrl.handleGetFavourites)
	api.GET("/profile", ctrl.handleGetMe)
	api.PATCH("/update_password", ctrl.handleChangePassword)
	api.DELETE("/delete_user", ctrl.handleDeleteUser)
	api.POST("/login", ctrl.handleLogin)
	api.GET("/fiilm/:id", ctrl.handleGetFilmByID)
	api.DELETE("/favorites/:id", ctrl.handleDeleteFromFavorites)
	api.POST("/accept_email", ctrl.handleAcceptEmail)
	api.POST("/comment/:id", ctrl.handleAddComment)
	api.DELETE("/comment/:id", ctrl.handleDeleteComment)
	api.POST("/favorites/:id", ctrl.handleAddToFavourites)
	api.GET("/genres", ctrl.handleGetAllGenres)
	api.GET("/:genre/:page", ctrl.handleGetFilmsByGenre)
	api.GET("/refresh", ctrl.handleRefresh)
	api.GET("/film/:name", ctrl.handleGetFilmsByName)
	api.GET("/comments/:id", ctrl.handleGetCommentsByFilmID)
}

// func to configure middlewares
func (ctrl *Controller) RegisterMiddlewares() {
	ctrl.logger.Info("registering middlewares")
	var middlewares = []echo.MiddlewareFunc{
		middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{echo.GET, echo.POST, echo.PATCH, echo.DELETE},
		}),
		middleware.Gzip(),
		middleware.Recover(),
		middleware.RequestIDWithConfig(
			middleware.RequestIDConfig{
				Skipper:      middleware.DefaultSkipper,
				Generator:    uuid.NewString,
				TargetHeader: echo.HeaderXRequestID,
			}),
		middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
			LogValuesFunc: ctrl.logValuesFunc,
			LogLatency:    true,
			LogRequestID:  true,
			LogMethod:     true,
			LogURI:        true,
		}),
	}
	ctrl.server.Use(middlewares...)
}

// func to improve logger
func (ctrl *Controller) logValuesFunc(_ echo.Context, v middleware.RequestLoggerValues) error {
	ctrl.logger.Info("Request",
		zap.String("uri", v.URI),
		zap.String("method", v.Method),
		zap.Duration("duration", v.Latency),
		zap.String("request-id", v.RequestID),
	)
	return nil
}
