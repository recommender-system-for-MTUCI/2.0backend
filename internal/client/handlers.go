package client

import (
	"encoding/json"
	"github.com/recommender-system-for-MTUCI/2.0backend/internal/models"
	"go.uber.org/zap"
	"strconv"
)

func (client *Client) HandleGetFilmID(id int) ([]int, error) {
	client.logger.Info("try to request film id")
	url := client.cfg.Clinet.Address() + strconv.Itoa(id)
	var data models.ResponseID
	resp, err := client.client.R().Get(url)
	if err != nil {
		client.logger.Error("Failed to perform request", zap.Error(err))
		return nil, err
	}
	err = json.Unmarshal(resp.Body(), &data)
	if err != nil {
		client.logger.Error("Failed to unmarshal response body", zap.Error(err))
		return nil, err
	}
	client.logger.Info("Successfully unmarshalled response", zap.Any("parsed_data", data))

	var filmIDs []int
	for _, film := range data.RecommendMovies {
		filmIDs = append(filmIDs, film)
	}
	client.logger.Info("Extracted film IDs", zap.Ints("filmIDs", filmIDs))

	return filmIDs, nil
}
