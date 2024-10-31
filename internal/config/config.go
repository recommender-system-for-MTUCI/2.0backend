package config

type Config struct {
	Server *Server
	JWT    *JWT
	SMTP   *SMTP
}

func New() (*Config, error) {
	cfg := &Config{
		Server: &Server{
			Host: "localhost",
			Port: 8080,
		},
		JWT: &JWT{
			AccessTime:  14555,
			RefreshTime: 53423,
			PublicKey:   "/home/relationskatie/backendForReccomenSystem/key/public.pem",
			PrivateKey:  "/home/relationskatie/backendForReccomenSystem/key/private.pem",
		},
		SMTP: &SMTP{
			SmtpServer:   "smtp.mail.ru",
			SmtpPort:     "587",
			From:         "reccomendsystem@mail.ru",
			PasswordPath: "/home/relationskatie/backendForReccomenSystem/key/password.txt",
		},
	}
	return cfg, nil
}
