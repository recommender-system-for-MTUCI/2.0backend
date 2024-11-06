package storage

import (
	"context"
	"github.com/google/uuid"
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
	const add = `INSERT INTO users(id, login, password, confirmation, code) VALUES ($1, $2, $3, $4, $5)`
	_, err := u.pgx.Exec(ctx, add, user.ID, user.Login, user.Password, user.Confirmation, user.Code)
	if err != nil {
		u.logger.Error("error adding user", zap.Error(err))
		return err
	}
	return nil
}

func (u *user) GetCodeFromDB(ctx context.Context, id uuid.UUID) (int, error) {
	const getCode = `SELECT code FROM users WHERE id = $1`
	var code int
	err := u.pgx.QueryRow(ctx, getCode, id).Scan(&code)
	if err != nil {
		u.logger.Error("error querying code", zap.Error(err))
		return 0, err
	}
	return code, err
}
func (u *user) GetPasswordFromDB(ctx context.Context, id uuid.UUID) (string, error) {
	const getPassword = `SELECT password FROM users WHERE id = $1`
	var password string
	err := u.pgx.QueryRow(ctx, getPassword, id).Scan(&password)
	if err != nil {
		u.logger.Error("error querying password", zap.Error(err))
		return "", err
	}
	return password, err

}
func (u *user) UpdatePassword(ctx context.Context, id uuid.UUID, newPassword string) error {
	const updatePassword = `UPDATE users SET password = $1 WHERE id = $2`
	_, err := u.pgx.Exec(ctx, updatePassword, newPassword, id)
	if err != nil {
		u.logger.Error("error updating password", zap.Error(err))
		return err
	}
	return nil
}
func (u *user) GetStatusFromUser(ctx context.Context, id uuid.UUID) (bool, error) {
	const getStatus = `SELECT confirmation FROM users WHERE id = $1`
	var confirmed bool
	err := u.pgx.QueryRow(ctx, getStatus, id).Scan(&confirmed)
	if err != nil {
		u.logger.Error("error querying status", zap.Error(err))
		return false, err
	}
	return confirmed, nil
}

func (u *user) DeleteUser(ctx context.Context, id uuid.UUID) error {
	const deleteUser = `DELETE FROM users WHERE id = $1`
	_, err := u.pgx.Exec(ctx, deleteUser, id)
	if err != nil {
		u.logger.Error("error deleting user", zap.Error(err))
		return err
	}
	return nil
}

func (u *user) GetUserIdByEmail(ctx context.Context, email string) (uuid.UUID, error) {
	const getUserIdByEmail = `SELECT id FROM users WHERE email = $1`
	var id uuid.UUID
	err := u.pgx.QueryRow(ctx, getUserIdByEmail, email).Scan(&id)
	if err != nil {
		u.logger.Error("error querying id by email", zap.Error(err))
		return uuid.Nil, err
	}
	return id, nil
}
