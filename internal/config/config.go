package config

type Config struct {
	Server *Server
}

func New() (*Config, error) {
	cfg := &Config{
		Server: &Server{
			Host: "localhost",
			Port: 8080,
		},
	}
	return cfg, nil
}
