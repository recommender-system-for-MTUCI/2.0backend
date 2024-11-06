package controller

import (
	"errors"
	"math/rand"
	"net/smtp"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
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
func validateRegistration(req *models.RequestRegister) error {
	if req.Login == "" || req.Password == "" {
		return errors.New("login or password is empty")
	}
	if len(req.Password) < 9 {
		return errors.New("password is too short")
	}
	return nil
}
func validatePassword(req *models.RequestChangePassword) error {
	if req.NewPassword == "" {
		return errors.New("new password is empty")
	}
	if len(req.NewPassword) < 9 {
		return errors.New("new password is too short")
	}
	return nil
}
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
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

/*func (ctrl *Controller) sendMessages(login string) error {
	errChan := make(chan error, 1)

	go func() {
		password, err := os.ReadFile(ctrl.cfg.SMTP.PasswordPath)
		if err != nil {
			ctrl.logger.Error("Failed to read SMTP password")
			errChan <- err
			return
		}

		from := ctrl.cfg.SMTP.From
		server := ctrl.cfg.SMTP.SmtpServer
		auth := authorization(from, string(password), server)
		smtpAddress := ctrl.cfg.SMTP.GetSmtpAddress()
		number := strconv.Itoa(generationRandomCode())
		code := []byte(number)

		maxRetries := 2
		for i := 0; i < maxRetries; i++ {
			err = smtp.SendMail(smtpAddress, auth, from, []string{login}, code)
			if err == nil {
				errChan <- nil
				return
			}

			ctrl.logger.Info("Failed to send message, retrying...")
			time.Sleep(time.Second)
		}

		ctrl.logger.Error("Failed to send message after retries")
		errChan <- err
	}()

	return <-errChan
}*/

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

func authorization(from string, password string, server string) smtp.Auth {
	auth := smtp.PlainAuth("", from, password, server)
	return auth
}

func (ctrl *Controller) generationRandomCode() int {
	generator := rand.New(rand.NewSource(time.Now().UnixNano()))
	code := generator.Intn(1000000)
	if code < 100000 {
		return code + 100000
	}
	return code
}

func decryptPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
