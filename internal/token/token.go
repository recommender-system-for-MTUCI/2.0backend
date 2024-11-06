package token

import (
	"crypto/rsa"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"github.com/recommender-system-for-MTUCI/2.0backend/internal/config"
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
	isAccess bool
}

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
func (t *Token) CreateToken(userID uuid.UUID, isAccess bool) (string, error) {
	var clock time.Duration
	if isAccess {
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
		isAccess: isAccess,
	}
	jwtToken, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(t.privateKey)
	if err != nil {
		return "", err
	}
	return jwtToken, nil
}

/*func (t *Token) ParseToken(tokenString string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenString, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {})
}*/
