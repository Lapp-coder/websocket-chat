package main

import (
	"context"
	"net/http"
)

type server struct {
	httpServer *http.Server
}

func newServer(host, port string) *server {
	return &server{
		httpServer: &http.Server{
			Addr: host + ":" + port,
		},
	}
}

func (s *server) start() error {
	return s.httpServer.ListenAndServe()
}

func (s *server) shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
