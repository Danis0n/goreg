package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/Danis0n/goreg/internal/goreg/httpprovider"
	"github.com/Danis0n/goreg/internal/goreg/server"
	"go.uber.org/zap"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	store       *ClientStore
	logger      *zap.Logger
	httpClient  HTTPClient
	registrator string
	errch       chan error
	closeCh     chan struct{}
	closeDoneCh chan struct{}
}

type RegisterRequest struct {
	Callback string `json:"callback"`
	Name     string `json:"name"`
	Port     int    `json:"port"`
}

type RegisterResponse struct {
	Hash string `json:"hash"`
}

const (
	maxRetries = 5
	callback   = "/callback"
)

func NewClient(cfg ClientConfig) (*Client, error) {
	stor, err := NewClientStore(cfg)
	if err != nil {
		return nil, err
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}

	return &Client{
		store:       stor,
		logger:      logger,
		registrator: cfg.Registrator,
		errch:       make(chan error),
		closeCh:     make(chan struct{}),
		closeDoneCh: make(chan struct{}),
		httpClient:  &http.Client{},
	}, nil
}

func NewClientWithStart(cfg ClientConfig) (*Client, error) {
	client, err := NewClient(cfg)
	if err != nil {
		return nil, err
	}

	client.Start()
	return client, nil
}

func (c *Client) Start() {
	defer close(c.closeDoneCh)

	c.StartListener(callback)
	c.doRegister()

	go func() {
		for {
			select {
			case <-c.closeCh:
				c.logger.Info("goreg->[client]: studtown")
				return
			case err := <-c.errch:
				c.logger.Error(err.Error())
			}
		}
	}()
}

func (c *Client) Stutdown() {
	c.doUnregister()
	close(c.closeCh)
	<-c.closeDoneCh
}

func (c *Client) StartListener(callback string) {
	http.HandleFunc(callback, func(w http.ResponseWriter, r *http.Request) {
		if err := server.ValidateHttpMethod(r.Method, http.MethodGet); err != nil {
			http.Error(w, err.Error(), http.StatusMethodNotAllowed)
			return
		}

		hash := r.URL.Query().Get("hash")
		if hash == "" {
			http.Error(w, "hash is required", http.StatusBadRequest)
			return
		}

		if err := c.Hash(hash); err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
		}

		w.WriteHeader(http.StatusOK)
	})
}

func (c *Client) Hash(hash string) error {
	if hash != c.store.Hash {
		return errors.New("goreg->[client]: hash dismatch")
	}
	return nil
}

func (g *Client) doRegister() {
	if g.store.Hash != "" {
		g.logger.Warn("goreg->[client]: already has hash")
		return
	}

	defer close(g.errch)

	b := &RegisterRequest{
		Callback: g.store.Callback,
		Name:     g.store.Name,
		Port:     g.store.Port,
	}

	reqBytes := new(bytes.Buffer)
	if err := json.NewEncoder(reqBytes).Encode(b); err != nil {
		g.logger.Error("goreg->[client]: request encoding error")
		g.errch <- err
		return
	}

	req, err := http.NewRequest(http.MethodPost, g.registrator, reqBytes)
	if err != nil {
		g.logger.Error("goreg->[client]: request create error")
		g.errch <- err
		return
	}

	for i := 0; i < maxRetries; i++ {
		data, err := httpprovider.Request(req, g.httpClient)
		if err != nil {
			g.logger.Error("goreg->[client]: request error")
			time.Sleep(time.Second)
			continue
		}

		var response RegisterResponse
		if err := json.Unmarshal(data, &response); err != nil {
			g.logger.Error("goreg->[client]: response unmarshal error")
			time.Sleep(time.Second)
			continue
		}

		g.store.Hash = response.Hash
		g.logger.Info("goreg->[client]: service was registered")
		return
	}

	g.logger.Error("goreg->[client]: registration failed after max retries")
}

func (g *Client) doUnregister() {
	url := g.registrator + "&name=" + g.store.Name + strconv.Itoa(g.store.Port)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		g.logger.Error("goreg->[client]: request create error")
		g.errch <- err
		return
	}

	for i := 0; i < maxRetries; i++ {
		_, err := httpprovider.Request(req, g.httpClient)
		if err != nil {
			g.logger.Error("goreg->[client]: request error")
			time.Sleep(time.Second)
			continue
		}

		g.logger.Info("goreg->[client]: service was unregistered")
		return
	}

	g.logger.Error("goreg->[client]: unregistration failed after max retries")
}
