package goreg

import (
	"github.com/Danis0n/goreg/internal/goreg/client"
	"github.com/Danis0n/goreg/internal/goreg/server"
)

func NewGoregClient(cfg client.ClientConfig) (*client.Client, error) {
	return client.NewClient(cfg)
}

func NewGoregClientWithStart(cfg client.ClientConfig) (*client.Client, error) {
	return client.NewClientWithStart(cfg)
}

func NewGoregServer(cfg server.ServerConfig) (*server.Server, error) {
	return server.NewServer(cfg)
}

func NewGoregServerWithStart(cfg server.ServerConfig) (*server.Server, error) {
	return server.NewServerWithStart(cfg)
}

func NewGoregClientConfig(registrator, callback, name string, port int) (client.ClientConfig, error) {
	return client.NewClientConfigWithName(registrator, callback, port, name)
}

func NewGoregClientConfigWithDefaults(registrator, callback string, port int) (client.ClientConfig, error) {
	return client.NewClientConfigWithDefaults(registrator, callback, port)
}

func NewGoregServerConfig(port int) (server.ServerConfig, error) {
	return server.NewServerConfig(port)
}
