package config

import "fmt"

type SMTP struct {
	SmtpServer   string
	SmtpPort     string
	From         string
	PasswordPath string
}

// func to return smtp addres
func (smtp SMTP) GetSmtpAddress() string {
	return fmt.Sprintf("%s:%s", smtp.SmtpServer, smtp.SmtpPort)
}
