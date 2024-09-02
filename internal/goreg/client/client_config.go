package client

import (
	"errors"

	"github.com/google/uuid"
)

type ClientConfig struct {
	Registrator string `yaml:"address"`
	Callback    string `yaml:"callback_address"`
	Name        string `yaml:"name"`
	Port        int    `yaml:"port"`
}

const (
	DefalutCallbackAddress = "callback"
)

func NewClientConfigWithDefaults(
	registratorAddress string,
	callbackAddress string,
	port int,
) (ClientConfig, error) {
	if err := validateClientSettings(registratorAddress, callbackAddress, port); err != nil {
		return ClientConfig{}, err
	}

	return ClientConfig{
		Registrator: registratorAddress,
		Callback:    callbackAddress,
		Name:        uuid.New().String(),
		Port:        port,
	}, nil
}

func NewClientConfigWithName(registrator string, callbackAddress string, port int, name string) (ClientConfig, error) {
	if err := validateClientSettings(registrator, callbackAddress, port); err != nil {
		return ClientConfig{}, err
	}

	return ClientConfig{
		Registrator: registrator,
		Callback:    callbackAddress,
		Name:        name,
		Port:        port,
	}, nil
}

func ValidateClientConfig(cfg ClientConfig) error {
	if err := validateClientSettings(cfg.Registrator, cfg.Callback, cfg.Port); err != nil {
		return err
	}
	return nil
}

func validateClientSettings(registrator string, callbackAddress string, port int) error {
	if callbackAddress == "" {
		return errors.New("callbackAddress invalid")
	}

	if registrator == "" {
		return errors.New("registrator invalid")
	}

	if port <= 0 {
		return errors.New("port invalid")
	}

	return nil
}
