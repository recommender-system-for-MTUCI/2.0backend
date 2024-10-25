package controller

import (
	"errors"
	"github.com/recommender-system-for-MTUCI/2.0backend/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func validateLogin(req *models.RequestLogin) error {
	if req.Login == "" || req.Password == "" {
		return errors.New("login or password is empty")
	}
	if len(req.Password) < 9 {
		return errors.New("password is too short")
	}
	return nil
}
func validateRegistration(req *models.RequestLogin) error {
	if req.Login == "" || req.Password == "" {
		return errors.New("login or password is empty")
	}
	if len(req.Password) < 9 {
		return errors.New("password is too short")
	}
	return nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 20)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
