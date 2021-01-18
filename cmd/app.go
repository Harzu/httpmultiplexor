package main

import (
	"context"
	"fmt"
	"httpmultiplexor/internal/config"
	"httpmultiplexor/internal/handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.InitConfig()

	h := handlers.New(cfg)
	h.Register()

	server := &http.Server{Addr: fmt.Sprintf(":%s", cfg.Port)}
	go func() {
		log.Printf("start server on %s port\n", cfg.Port)
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("failed to start server %s\n", err.Error())
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	<-signalChan

	log.Println("server shutdown")
	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatalf("failed to shutdown server: %s", err.Error())
	}
}