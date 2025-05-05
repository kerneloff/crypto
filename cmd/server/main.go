package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/kerneloff/crypto/internal/server"
)

func main() {
	srv := server.NewServer()

	// Запускаем сервер в отдельной горутине
	go func() {
		if err := srv.Run(":8080"); err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	// Ожидаем сигнал для graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	// Здесь можно добавить graceful shutdown
} 