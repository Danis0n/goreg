package client

import (
	"testing"
)

func TestNewClientStore(t *testing.T) {
	cfg := ClientConfig{
		Registrator: "http://registrator.url",
		Callback:    "http://callback.url",
		Name:        "test-client",
		Port:        8080,
	}

	store, err := NewClientStore(cfg)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if store.Callback != cfg.Callback {
		t.Errorf("Expected Callback %s, got %s", cfg.Callback, store.Callback)
	}

	if store.Name != cfg.Name {
		t.Errorf("Expected Name %s, got %s", cfg.Name, store.Name)
	}

	if store.Port != cfg.Port {
		t.Errorf("Expected Port %d, got %d", cfg.Port, store.Port)
	}

	if store.Hash != "" {
		t.Errorf("Expected Hash to be empty, got %s", store.Hash)
	}

	if store.logger == nil {
		t.Error("Expected logger to be initialized, got nil")
	}
}

func TestNewClientStore_InvalidConfig(t *testing.T) {
	invalidCfg := ClientConfig{
		Registrator: "",
		Callback:    "",
		Name:        "invalid-client",
		Port:        -1,
	}

	_, err := NewClientStore(invalidCfg)
	if err == nil {
		t.Fatal("Expected error for invalid config, got nil")
	}
}
