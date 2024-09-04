package main

import "github.com/Danis0n/goreg/internal/goreg"

func setupServer() {

	cfg, err := goreg.NewGoregServerConfig(8079)
	if err != nil {
		// do smth
	}

	server, err := goreg.NewGoregServer(cfg)
	if err != nil {
		// do smth
	}

	server.Start()
}

func setupClient() {

	cfg, err := goreg.NewGoregClientConfig(
		"https://some-api.com",
		"https://some-callback.com",
		"test",
		90,
	)
	if err != nil {
		// do smth
	}

	client, err := goreg.NewGoregClient(cfg)
	if err != nil {
		// do smth
	}

	client.Start()
}

func main() {
	setupServer()
	setupClient()
}
