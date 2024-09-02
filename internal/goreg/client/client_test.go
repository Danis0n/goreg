package client

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

func TestClientDoRegister_Success(t *testing.T) {
	cfg := ClientConfig{
		Registrator: "http://registrator.url",
		Callback:    "http://callback.url",
		Name:        "test-client",
		Port:        8080,
	}

	client, err := NewClient(cfg)
	assert.NoError(t, err)

	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			response := RegisterResponse{Hash: "test-hash"}
			respBytes, _ := json.Marshal(response)

			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBuffer(respBytes)),
			}, nil
		},
	}

	client.httpClient = mockClient
	client.doRegister()

	assert.Equal(t, "test-hash", client.store.Hash)
}

func TestClientDoRegister_Failure(t *testing.T) {
	cfg := ClientConfig{
		Registrator: "http://registrator.url",
		Callback:    "http://callback.url",
		Name:        "test-client",
		Port:        8080,
	}

	client, err := NewClient(cfg)
	assert.NoError(t, err)

	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       io.NopCloser(bytes.NewBufferString("")),
			}, nil
		},
	}

	client.httpClient = mockClient
	client.doRegister()

	assert.Empty(t, client.store.Hash)
}
