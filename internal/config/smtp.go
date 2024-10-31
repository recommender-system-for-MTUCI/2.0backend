package config

type SMTP struct {
	SmtpServer   string
	SmtpPort     string
	From         string
	PasswordPath string
}

func (smtp SMTP) GetSmtpAddress() string {
	return smtp.SmtpServer + ":" + smtp.SmtpPort
}
