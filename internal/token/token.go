package token

import (
	"crypto/rsa"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"github.com/recommender-system-for-MTUCI/2.0backend/internal/config"
	"github.com/recommender-system-for-MTUCI/2.0backend/internal/models"
	"os"
	"time"
)

var _ JWT = (*Token)(nil)

type Token struct {
	publicKey   *rsa.PublicKey
	privateKey  *rsa.PrivateKey
	accessTime  int
	refreshTime int
}
type tokenClaims struct {
	jwt.RegisteredClaims
	IsAccess bool `json:"isAccess"`
}

// create token config
func NewToken(cfg *config.Config) (*Token, error) {
	publicKeyPath, err := os.ReadFile(cfg.JWT.PublicKey)
	if err != nil {
		log.Error("got error to open file")
		return nil, err
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyPath)
	if err != nil {
		return nil, err
	}
	privateKeyPath, err := os.ReadFile(cfg.JWT.PrivateKey)
	if err != nil {
		return nil, err
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyPath)
	if err != nil {
		return nil, err
	}
	token := &Token{
		publicKey:   publicKey,
		privateKey:  privateKey,
		accessTime:  cfg.JWT.AccessTime,
		refreshTime: cfg.JWT.RefreshTime,
	}
	return token, nil
}

// create new token with claims
func (t *Token) CreateToken(userID uuid.UUID, isAccess bool) (string, error) {
	var clock time.Duration
	if isAccess == true {
		clock = time.Duration(t.accessTime) * time.Minute
	} else {
		clock = time.Duration(t.refreshTime) * time.Minute
	}
	claims := tokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    userID.String(),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(clock)),
		},
		IsAccess: isAccess,
	}
	jwtToken, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(t.privateKey)
	if err != nil {
		return "", err
	}
	return jwtToken, nil
}

// parse token claims
func (t *Token) ParseToken(tokenString string) (*models.UserData, error) {
	token, err := jwt.ParseWithClaims(tokenString, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			log.Error("unexpected signing method")
			return nil, errors.New("unexpected signing method")
		}
		return t.publicKey, nil
	})
	if err != nil {
		log.Error("Failed to parse token:", err)
		return nil, err
	}
	claims, ok := token.Claims.(*tokenClaims)
	if !ok || !token.Valid {
		log.Error("token is not ok or not valid")
		return nil, errors.New("token is not ok or not valid")
	}
	if claims.ExpiresAt.Time.Before(time.Now()) {
		log.Error("token has expired")
		return nil, errors.New("token has expired")
	}
	data := &models.UserData{
		ID:       uuid.MustParse(claims.Issuer),
		IsAccess: claims.IsAccess,
	}
	return data, nil
}
