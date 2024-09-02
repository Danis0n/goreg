package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"

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

	req, err := http.NewRequest("POST", g.registrator, reqBytes)
	if err != nil {
		g.logger.Error("goreg->[client]: request create error")
		g.errch <- err
		return
	}

	for i := 0; i < maxRetries; i++ {
		data, err := g.request(req)
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
	req, err := http.NewRequest(
		"DELETE",
		g.registrator+"&name="+g.store.Name+strconv.Itoa(g.store.Port),
		nil)
	if err != nil {
		g.logger.Error("goreg->[client]: request create error")
		g.errch <- err
		return
	}

	for i := 0; i < maxRetries; i++ {
		_, err := g.request(req)
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

func (c *Client) request(req *http.Request) ([]byte, error) {
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("goreg: bad status code: " + res.Status)
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return bodyBytes, nil
}
