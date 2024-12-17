package config

import (
	"fmt"
	"github.com/labstack/gommon/log"
	"os"
)

type Postgres struct {
	Host         string
	Port         string
	User         string
	PasswordPath string
	Database     string
}

// func to return postgressql address
func (p Postgres) GetAddressPostgres() string {
	file, err := os.ReadFile(p.PasswordPath)
	if err != nil {
		log.Error(err)
	}
	password := string(file)
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", p.User, password, p.Host, p.Port, p.Database)
}
