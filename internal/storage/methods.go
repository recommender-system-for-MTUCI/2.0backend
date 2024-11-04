package storage

import (
	"context"
	"github.com/recommender-system-for-MTUCI/2.0backend/internal/models"
)

type UserStorage interface {
	AddUserInDB(ctx context.Context, user *models.DTORegister) error
}

type DB interface {
	User() UserStorage
}
