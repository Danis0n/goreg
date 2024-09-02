package server

import (
	"errors"
)

type ServerConfig struct {
	Port int `yaml:"port"`
}

func NewServerConfig(port int) (ServerConfig, error) {
	if err := validateServerSettings(port); err != nil {
		return ServerConfig{}, err
	}

	return ServerConfig{
		Port: port,
	}, nil
}

func ValidateServerConfig(cfg ServerConfig) error {
	if err := validateServerSettings(cfg.Port); err != nil {
		return err
	}
	return nil
}

func validateServerSettings(port int) error {
	if port <= 0 {
		return errors.New("port invalid")
	}

	return nil
}
