package server

import (
	"errors"
	"sync"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Service struct {
	Name     string
	Hash     string
	Callback string
}

type ServerStore struct {
	logger   *zap.Logger
	rwmu     *sync.RWMutex
	services map[string]*Service
}

func NewServerStore(logger *zap.Logger) (*ServerStore, error) {
	if logger == nil {
		return nil, errors.New("logger invalid")
	}

	return &ServerStore{
		logger:   logger,
		rwmu:     &sync.RWMutex{},
		services: make(map[string]*Service),
	}, nil
}

func (g *ServerStore) Get(name string) (*Service, error) {
	g.rwmu.RLock()
	defer g.rwmu.RUnlock()

	services, ok := g.services[name]
	if !ok {
		return nil, errors.New("registrator [server]: service not found")
	}

	return services, nil
}

func (g *ServerStore) Set(name string, callback string) error {
	g.rwmu.Lock()
	defer g.rwmu.Unlock()

	_, ok := g.services[name]
	if ok {
		return errors.New("registrator [server]: server already exists")
	}

	g.services[name] = &Service{
		Name:     name,
		Hash:     uuid.New().String(),
		Callback: callback,
	}
	g.logger.Info("Registrator [server]: service: " + name + " was registered")

	return nil
}

func (g *ServerStore) GetAll() []*Service {
	g.rwmu.RLock()
	defer g.rwmu.RUnlock()

	servers := make([]*Service, 0, len(g.services))
	for _, value := range g.services {
		servers = append(servers, value)
	}

	return servers
}

func (g *ServerStore) Delete(key string) error {
	g.rwmu.Lock()
	defer g.rwmu.Unlock()

	_, ok := g.services[key]
	if !ok {
		return errors.New("Registrator [server]: key{" + key + "} doesn't exists")
	}

	delete(g.services, key)
	g.logger.Info("Registrator [server]: service: {" + key + "} was removed")

	return nil
}
