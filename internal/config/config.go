package config

type Config struct {
	Server   *Server
	JWT      *JWT
	SMTP     *SMTP
	Postgres *Postgres
}

// func to create new instance config
func New() (*Config, error) {
	cfg := &Config{
		Server: &Server{
			Host: "localhost",
			Port: 8080,
		},
		JWT: &JWT{
			AccessTime:  15,
			RefreshTime: 90,
			PublicKey:   "/home/relationskatie/backendForReccomenSystem/key/public.pem",
			PrivateKey:  "/home/relationskatie/backendForReccomenSystem/key/private.pem",
		},
		SMTP: &SMTP{
			SmtpServer:   "smtp.mail.ru",
			SmtpPort:     "587",
			From:         "reccomendsystem@mail.ru",
			PasswordPath: "/home/relationskatie/backendForReccomenSystem/key/password.txt",
		},
		Postgres: &Postgres{
			Host:         "localhost",
			Port:         "5432",
			User:         "relationskatie",
			PasswordPath: "/home/relationskatie/backendForReccomenSystem/key/passwordPostgres.txt",
			Database:     "recommend_system",
		},
	}
	return cfg, nil
}
