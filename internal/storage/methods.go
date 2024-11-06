package storage

import (
	"context"
	"github.com/google/uuid"
	"github.com/recommender-system-for-MTUCI/2.0backend/internal/models"
)

type UserStorage interface {
	AddUserInDB(ctx context.Context, user *models.DTORegister) error
	GetCodeFromDB(ctx context.Context, id uuid.UUID) (int, error)
	GetPasswordFromDB(ctx context.Context, id uuid.UUID) (string, error)
	UpdatePassword(ctx context.Context, id uuid.UUID, newPassword string) error
	GetStatusFromUser(ctx context.Context, id uuid.UUID) (bool, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
	GetUserIdByEmail(ctx context.Context, email string) (uuid.UUID, error)
}

type DB interface {
	User() UserStorage
}
