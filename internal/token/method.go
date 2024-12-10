package token

import (
	"github.com/google/uuid"
	"github.com/recommender-system-for-MTUCI/2.0backend/internal/models"
)

type JWT interface {
	CreateToken(userID uuid.UUID, isAccess bool) (string, error)
	ParseToken(tokenString string) (*models.UserData, bool, error)
}
