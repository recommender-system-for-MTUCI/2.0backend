package controller

import (
	"errors"
	"github.com/google/uuid"
	"github.com/recommender-system-for-MTUCI/2.0backend/internal/models"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"net/smtp"
	"os"
	"strconv"
	"time"
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
func validateRegistration(req models.RequestRegister) error {
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

func (ctrl *Controller) sendMessages(login string) error {
	password, err := os.ReadFile(ctrl.cfg.SMTP.PasswordPath)
	if err != nil {
		return err
	}
	from := ctrl.cfg.SMTP.From
	server := ctrl.cfg.SMTP.SmtpServer
	auth := authorization(from, string(password), server)
	smtpAddress := ctrl.cfg.SMTP.GetSmtpAddress()
	number := strconv.Itoa(generationRandomCode())
	code := []byte(number)
	err = smtp.SendMail(smtpAddress, auth, from, []string{login}, code)
	if err != nil {
		ctrl.logger.Info("Failed to send message")
	}
	return err
}

func authorization(from string, password string, server string) smtp.Auth {
	auth := smtp.PlainAuth("", from, password, server)
	return auth
}

func generationRandomCode() int {
	generator := rand.New(rand.NewSource(time.Now().UnixNano()))
	code := generator.Intn(1000000)
	return code
}
