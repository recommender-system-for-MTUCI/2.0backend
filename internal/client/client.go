package client

import (
	"github.com/go-resty/resty/v2"
	"github.com/recommender-system-for-MTUCI/2.0backend/internal/config"
	"go.uber.org/zap"
)

type Client struct {
	logger *zap.Logger
	cfg    *config.Config
	client *resty.Client
}

func New(logger *zap.Logger, cfg *config.Config) (*Client, error) {
	client := &Client{
		logger: logger,
		cfg:    cfg,
		client: resty.New(),
	}
	return client, nil

}
