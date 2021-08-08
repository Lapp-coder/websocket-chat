package main

import (
	"context"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	hub := newHub()
	handler := newHandler(hub)
	handler.initRoutes()
	go hub.run()

	host := os.Getenv("CHAT_HOST")
	port := os.Getenv("CHAT_PORT")
	srv := newServer(host + ":" + port)

	go func() {
		if err := srv.start(); err != nil {
			logrus.Errorf("failed to start server: %s", err.Error())
		}
	}()

	logrus.Info("server started")

	shutdown := make(chan os.Signal)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown

	logrus.Info("server shutdown")

	if err := srv.shutdown(context.Background()); err != nil {
		logrus.Errorf("failed to graceful shutdown server: %s", err.Error())
	}
}
