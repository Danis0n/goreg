package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/Danis0n/goreg/internal/goreg/httpprovider"
	"go.uber.org/zap"
)

type Server struct {
	logger      *zap.Logger
	store       *ServerStore
	errch       chan error
	closeCh     chan struct{}
	closeDoneCh chan struct{}
	port        int
	httpClient  httpprovider.HttpClient
}

func NewServer(cfg ServerConfig) (*Server, error) {
	if err := ValidateServerConfig(cfg); err != nil {
		return nil, err
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}

	stor, err := NewServerStore(logger)
	if err != nil {
		return nil, err
	}

	return &Server{
		logger:      logger,
		store:       stor,
		errch:       make(chan error),
		closeCh:     make(chan struct{}),
		closeDoneCh: make(chan struct{}),
		port:        cfg.Port,
	}, nil
}

func NewServerWithStart(cfg ServerConfig) (*Server, error) {
	registrator, err := NewServer(cfg)
	if err != nil {
		return nil, err
	}

	registrator.Start()
	return registrator, nil
}

func (g *Server) Start() {
	defer close(g.closeDoneCh)

	g.startServer(g.port)

	go func() {
		for {
			select {
			case <-g.closeCh:
				g.logger.Info("goreg->[server]: shutdown")
				return
			case err := <-g.errch:
				g.logger.Error(err.Error())
			case <-time.After(time.Minute):
				g.checkServicesAvailability()
			}
		}
	}()
}

func (g *Server) startServer(port int) error {
	http.HandleFunc("/set", g.SetHandler)
	http.HandleFunc("/delete", g.DeleteHandler)
	http.HandleFunc("/getall", g.GetAllHandler)
	http.HandleFunc("/get", g.GetHandler)

	g.logger.Info("Server was started at port: " + strconv.Itoa(port))
	return http.ListenAndServe(":"+strconv.Itoa(port), nil)
}

func (g *Server) checkServicesAvailability() {
	for _, service := range g.store.services {
		go func() {
			g.logger.Info("goreg->[server]: check service availability: " + service.Name)
			g.checkServiceAvailability(*service)
		}()
	}
}

func (g *Server) checkServiceAvailability(service Service) {
	url := service.Callback + "?hash= " + service.Hash

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		g.logger.Error("goreg->[client]: request create error")
		g.errch <- err
		return
	}

	_, err = httpprovider.Request(req, g.httpClient)
	if err != nil {
		g.logger.Error("goreg->[client]: request send error")
		g.errch <- err
		return
	}
}

func (g *Server) GetHandler(w http.ResponseWriter, r *http.Request) {
	if err := ValidateHttpMethod(r.Method, http.MethodGet); err != nil {
		http.Error(w, err.Error(), http.StatusMethodNotAllowed)
		return
	}

	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	service, err := g.store.Get(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(service)
}

func (g *Server) SetHandler(w http.ResponseWriter, r *http.Request) {
	if err := ValidateHttpMethod(r.Method, http.MethodPost); err != nil {
		http.Error(w, err.Error(), http.StatusMethodNotAllowed)
		return
	}

	var svc Service
	if err := json.NewDecoder(r.Body).Decode(&svc); err != nil {
		g.logger.Error("invalid input")
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	if svc.Name == "" || svc.Callback == "" {
		g.logger.Error("name and callback are required")
		http.Error(w, "name and callback are required", http.StatusBadRequest)
		return
	}

	if err := g.store.Set(svc.Name, svc.Callback); err != nil {
		g.logger.Error("failed to set service: " + err.Error())
		http.Error(w, "failed to set service: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (g *Server) GetAllHandler(w http.ResponseWriter, r *http.Request) {
	if err := ValidateHttpMethod(r.Method, http.MethodGet); err != nil {
		http.Error(w, err.Error(), http.StatusMethodNotAllowed)
		return
	}

	services := g.store.GetAll()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(services)
}

func (g *Server) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	if err := ValidateHttpMethod(r.Method, http.MethodDelete); err != nil {
		http.Error(w, err.Error(), http.StatusMethodNotAllowed)
		return
	}

	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	if err := g.store.Delete(name); err != nil {
		http.Error(w, "Failed to delete service", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func ValidateHttpMethod(method string, requiredMethod string) error {
	if method != requiredMethod {
		return errors.New("method not allowed")
	}
	return nil
}
