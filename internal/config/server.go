package config

import "strconv"

type Server struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

func (s Server) GetAddress() string {
	return s.Host + ":" + strconv.Itoa(s.Port)
}
