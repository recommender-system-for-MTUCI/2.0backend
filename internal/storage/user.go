package storage

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/recommender-system-for-MTUCI/2.0backend/internal/models"
	"go.uber.org/zap"
)

var _ UserStorage = (*user)(nil)

type user struct {
	logger *zap.Logger
	pgx    *pgxpool.Pool
}

func NewUser(logger *zap.Logger, pgxPool *pgxpool.Pool) (*user, error) {
	user := &user{
		logger: logger,
		pgx:    pgxPool,
	}
	return user, nil
}

func (u *user) AddUserInDB(ctx context.Context, user *models.DTORegister) error {
	const add = `INSERT INTO users(id, login, password, confirmation) VALUES ($1, $2, $3, $4)`
	_, err := u.pgx.Exec(ctx, add, user.ID, user.Login, user.Password, user.Confirmation)
	if err != nil {
		u.logger.Error("error adding user", zap.Error(err))
		return err
	}
	return nil
}
