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

// create new instanse user storage
func NewUser(logger *zap.Logger, pgxPool *pgxpool.Pool) (*user, error) {
	user := &user{
		logger: logger,
		pgx:    pgxPool,
	}
	logger.Info("successfully created user storage")
	return user, nil
}

// func to add user in db
func (u *user) AddUserInDB(ctx context.Context, user *models.DTORegister) error {
	const add = `INSERT INTO users(id, login, password, is_active, code) VALUES ($1, $2, $3, $4, $5)`
	_, err := u.pgx.Exec(ctx, add, user.ID, user.Login, user.Password, user.Confirmation, user.Code)
	if err != nil {
		u.logger.Error("error adding user", zap.Error(err))
		return err
	}
	u.logger.Info("successfully added user in db")
	return nil
}

// func get code to accept email from db
func (u *user) GetCodeFromDB(ctx context.Context, id uuid.UUID) (int, error) {
	const getCode = `SELECT code FROM users WHERE id = $1`
	var code int
	err := u.pgx.QueryRow(ctx, getCode, id).Scan(&code)
	if err != nil {
		u.logger.Error("error querying code", zap.Error(err))
		return 0, err
	}
	u.logger.Info("successfully retrieved code")
	return code, err
}

// func to get hash password from db
func (u *user) GetPasswordFromDB(ctx context.Context, id uuid.UUID) (string, error) {
	const getPassword = `SELECT password FROM users WHERE id = $1`
	var password string
	err := u.pgx.QueryRow(ctx, getPassword, id).Scan(&password)
	if err != nil {
		u.logger.Error("error querying password", zap.Error(err))
		return "", err
	}
	u.logger.Info("successfully retrieved password")
	return password, err

}

// func to update user password
func (u *user) UpdatePassword(ctx context.Context, id uuid.UUID, newPassword string) error {
	const updatePassword = `UPDATE users SET password = $1 WHERE id = $2`
	_, err := u.pgx.Exec(ctx, updatePassword, newPassword, id)
	if err != nil {
		u.logger.Error("error updating password", zap.Error(err))
		return err
	}
	u.logger.Info("successfully updated password")
	return nil
}

// func to get user status from db, that means user aceepted email or not
func (u *user) GetStatusFromUser(ctx context.Context, id uuid.UUID) (bool, error) {
	const getStatus = `SELECT is_active FROM users WHERE id = $1`
	var confirmed bool
	err := u.pgx.QueryRow(ctx, getStatus, id).Scan(&confirmed)
	if err != nil {
		u.logger.Error("error querying status", zap.Error(err))
		return false, err
	}
	u.logger.Info("successfully retrieved status")
	return confirmed, nil
}

// func to delet user from db
func (u *user) DeleteUser(ctx context.Context, id uuid.UUID) error {
	const deleteUser = `DELETE FROM users WHERE id = $1`
	_, err := u.pgx.Exec(ctx, deleteUser, id)
	if err != nil {
		u.logger.Error("error deleting user", zap.Error(err))
		return err
	}
	u.logger.Info("successfully deleted user")
	return nil
}

// func to get user ID by email
func (u *user) GetUserIdByEmail(ctx context.Context, email string) (uuid.UUID, string, error) {
	const getUserIdByEmail = `SELECT id, password FROM users WHERE login = $1`
	var id uuid.UUID
	var password string
	err := u.pgx.QueryRow(ctx, getUserIdByEmail, email).Scan(&id, &password)
	if err != nil {
		u.logger.Error("error querying id by email", zap.Error(err))
		return uuid.Nil, "", err
	}
	u.logger.Info("successfully retrieved user by email")
	return id, password, nil
}

// func to update user status
func (u *user) UpdateUserStatus(ctx context.Context, id uuid.UUID) error {
	const updateStatus = `UPDATE users SET is_active = $1 WHERE id = $2`
	_, err := u.pgx.Exec(ctx, updateStatus, true, id)
	if err != nil {
		u.logger.Error("error updating status", zap.Error(err))
		return err
	}
	u.logger.Info("successfully updated status")
	return nil

}

// func to return user account
func (u *user) GetMe(ctx context.Context, id uuid.UUID) (string, error) {
	const getMe = `SELECT login FROM users WHERE id = $1`
	var login string
	err := u.pgx.QueryRow(ctx, getMe, id).Scan(&login)
	if err != nil {
		u.logger.Error("error querying user", zap.Error(err))
		return "", err
	}
	u.logger.Info("successfully retrieved me")
	return login, nil
}
