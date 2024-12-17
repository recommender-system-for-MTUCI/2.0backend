package controller

import (
	"errors"
	"github.com/labstack/gommon/log"
	"go.uber.org/zap"
	"math/rand"
	"net/http"
	"net/smtp"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/recommender-system-for-MTUCI/2.0backend/internal/models"
	"golang.org/x/crypto/bcrypt"
)

// func to validate email and password for login
func validateLogin(req models.RequestLogin) error {
	if req.Login == "" || req.Password == "" {
		return errors.New("login or password is empty")
	}
	if len(req.Password) < 9 {
		return errors.New("password is too short")
	}
	return nil
}

// func to validate email and password for registration
func validateRegistration(req models.RequestRegister) error {
	if req.Login == "" || req.Password == "" {
		return errors.New("login or password is empty")
	}
	if len(req.Password) < 9 {
		return errors.New("password is too short")
	}
	return nil
}

// func to validate password
func validatePassword(req models.RequestChangePassword) error {
	if req.NewPassword == "" {
		return errors.New("new password is empty")
	}
	if len(req.NewPassword) < 9 {
		return errors.New("new password is too short")
	}
	return nil
}

// func to hash password
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// func to encrypt password
func encryptPassword(password string, hashPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
	if err != nil {
		return err
	}
	return err
}

// func to generate access and refresh tokens
func (ctrl *Controller) generateAccessAndRefreshToken(userID uuid.UUID) (accessToken, refreshToken string, err error) {
	accessToken, err = ctrl.token.CreateToken(userID, true)
	if err != nil {
		return
	}
	refreshToken, err = ctrl.token.CreateToken(userID, false)
	if err != nil {
		return
	}
	return accessToken, refreshToken, nil
}

// func to accept user email
func (ctrl *Controller) sendMessages(login string, number int) error {
	password, err := os.ReadFile(ctrl.cfg.SMTP.PasswordPath)
	if err != nil {
		ctrl.logger.Error("Failed to read SMTP password")
		return err
	}
	from := ctrl.cfg.SMTP.From
	server := ctrl.cfg.SMTP.SmtpServer
	auth := authorization(from, string(password), server)
	smtpAddress := ctrl.cfg.SMTP.GetSmtpAddress()
	num := strconv.Itoa(number)
	code := []byte(num)
	maxRetries := 2
	for i := 0; i < maxRetries; i++ {
		err = smtp.SendMail(smtpAddress, auth, from, []string{login}, code)
		if err == nil {
			return nil
		}

		ctrl.logger.Info("Failed to send message, retrying...")
		time.Sleep(time.Second)
	}

	ctrl.logger.Error("Failed to send message after retries")
	return err
}

// func to authorization on smtp client
func authorization(from string, password string, server string) smtp.Auth {
	auth := smtp.PlainAuth("", from, password, server)
	return auth
}

// func to generation random code for accept email
func generationRandomCode() int {
	generator := rand.New(rand.NewSource(time.Now().UnixNano()))
	code := generator.Intn(1000000)
	if code < 100000 {
		return code + 100000
	}
	return code
}

// func to get jwt token from header
func tokenInHeader(req *http.Request) (string, error) {
	authHeader := req.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("Authorization header is empty")
	}
	token := strings.SplitN(authHeader, "Bearer ", 2)
	if len(token) != 2 {
		log.Error("Authorization header format is invalid")
		return "", errors.New("Authorization header is invalid")
	}
	if token[1] == "" {
		log.Error("token is empty")
		return "", errors.New("Authorization header is empty")
	}
	return strings.TrimSpace(token[1]), nil
}

// func to get userID from jwt token
func (ctrl *Controller) getUserId(req *http.Request) (uuid.UUID, error) {
	token, err := tokenInHeader(req)
	if err != nil {
		ctrl.logger.Error("Failed to get token")
		ctrl.logger.Info(token)
		return uuid.Nil, err
	}
	data, err := ctrl.token.ParseToken(token)
	if err != nil {
		ctrl.logger.Error("Failed to parse token")
		return uuid.Nil, err
	}
	if data.IsAccess != true {
		ctrl.logger.Error("Invalid token")
		return uuid.Nil, errors.New("Invalid token")
	}
	id := data.ID
	log.Info(id)
	return id, nil

}

// func to validate Refresh token
func (ctrl *Controller) validateRefreshToken(refreshToken string) (uuid.UUID, error) {
	data, err := ctrl.token.ParseToken(refreshToken)
	if err != nil {
		ctrl.logger.Error("failed to parse refresh token", zap.Error(err))
		return uuid.Nil, err
	}
	if data.IsAccess != false {
		ctrl.logger.Error("provided token is not a refresh token")
		return uuid.Nil, errors.New("invalid token type: expected refresh token")
	}
	ctrl.logger.Info("refresh token successfully validated")
	return data.ID, nil
}
