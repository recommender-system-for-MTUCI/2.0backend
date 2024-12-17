package storage

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

var _ DB = (*Storage)(nil)

type Storage struct {
	logger *zap.Logger
	pgx    *pgxpool.Pool
	user   *user
	film   *film
}

// func to create neew instanse storage
func NewStorage(logger *zap.Logger, pgx *pgxpool.Pool) (*Storage, error) {
	user, err := NewUser(logger, pgx)
	if err != nil {
		logger.Error("failed to create user storage", zap.Error(err))
		return nil, err
	}
	film, err := NewFilm(logger, pgx)
	if err != nil {
		logger.Error("failed to create film storage", zap.Error(err))
		return nil, err
	}
	storage := &Storage{
		logger: logger,
		pgx:    pgx,
		user:   user,
		film:   film,
	}
	logger.Info("successfully created storage")
	return storage, nil

}
func (s *Storage) User() UserStorage {
	return s.user

}
func (s *Storage) Film() FilmStorage {
	return s.film
}
