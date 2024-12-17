package config

import (
	"fmt"
)

type Server struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

// func to return server addres
func (s Server) GetAddress() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}
