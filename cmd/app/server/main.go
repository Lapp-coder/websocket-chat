package main

import (
	"context"
	"github.com/Lapp-coder/websocket-chat/internal/app/server"
	"github.com/sirupsen/logrus"
	"net/rpc"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	hub := server.NewHub()
	handler := server.NewHandler(hub)
	rpc.Register(handler)
	handler.InitRoutes()
	go hub.Listen()

	host := os.Getenv("CHAT_HOST")
	port := os.Getenv("CHAT_PORT")
	if host == "" || port == "" {
		logrus.Fatalf("incorrect address to start the server")
	}

	// Запуск сервера в go-рутине для его плавной остановки
	srv := server.NewServer(host + ":" + port)
	go func() {
		if err := srv.Start(); err != nil {
			logrus.Errorf("failed to start server: %s", err.Error())
		}
	}()

	logrus.Info("server started")

	// Ожидание на получение одного из системных сигналов (SIGINT, SIGTERM) для продолжение выполнения функции main
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown

	logrus.Info("server shutdown")

	// Плавная остановка сервера
	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("failed to graceful shutdown server: %s", err.Error())
	}
}
