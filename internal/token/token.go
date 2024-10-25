package token

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/recommender-system-for-MTUCI/2.0backend/internal/models"
	"time"
)

func NewToken(user models.DTOLogin, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodRS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["uuid"] = user.ID
	claims["exp"] = time.Now().Add(duration).Unix()
	return "", nil
}
