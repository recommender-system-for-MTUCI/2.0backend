package config

import (
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

func (p Postgres) GetAddressPostgres() string {
	file, err := os.ReadFile(p.PasswordPath)
	if err != nil {
		log.Panic(err)
	}
	password := string(file)
	return "postgres" + "://" + p.User + ":" + password + "@" + p.Host + ":" + p.Port + "/" + p.Database + "?sslmode=disable"
}
