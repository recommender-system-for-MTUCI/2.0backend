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

func NewStorage(logger *zap.Logger, pgx *pgxpool.Pool) (*Storage, error) {
	user, err := NewUser(logger, pgx)
	if err != nil {
		return nil, err
	}
	film, err := NewFilm(logger, pgx)
	if err != nil {
		return nil, err
	}
	storage := &Storage{
		logger: logger,
		pgx:    pgx,
		user:   user,
		film:   film,
	}
	return storage, nil

}
func (s *Storage) User() UserStorage {
	return s.user

}
func (s *Storage) Film() FilmStorage {
	return s.film
}
