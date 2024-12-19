package config

import "fmt"

type Clinet struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

func (c Clinet) Address() string {
	return fmt.Sprintf("http://%s:%d/predict/", c.Host, c.Port)
}
