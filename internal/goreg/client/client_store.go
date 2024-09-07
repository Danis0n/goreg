package client

import (
	"go.uber.org/zap"
)

type ClientStore struct {
	logger   *zap.Logger
	Hash     string
	Callback string
	Name     string
	Port     int
}

func NewClientStore(cfg ClientConfig) (*ClientStore, error) {
	if err := ValidateClientConfig(cfg); err != nil {
		return nil, err
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}

	return &ClientStore{
		logger:   logger,
		Callback: cfg.Callback + callback,
		Name:     cfg.Name,
		Port:     cfg.Port,
		Hash:     "",
	}, nil
}
