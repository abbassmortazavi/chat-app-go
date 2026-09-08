package main

import (
	"backend/internal/config"
	"backend/internal/db"
	"backend/internal/middlewares"
	"backend/internal/routes"
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := config.LoadConfig()
	db.Init(cfg.DBPath, cfg.DBName)
	defer db.CloseDB()

	mux := routes.RegisterRoutes()
	middlewareMux := middlewares.LoggingMiddleware(mux)

	server := &http.Server{
		Addr:    cfg.HttpServer.HttpAddress,
		Handler: middlewareMux,
	}

	serverErrCh := make(chan error, 1)

	go func() {
		log.Println("Starting server on port 8084")
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrCh <- err
			return
		}
		serverErrCh <- nil
	}()

	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-shutdownCh:
		log.Printf("Received signal %q, shutting down", sig)
		shutdown(server)
	case err := <-serverErrCh:
		if err != nil {
			log.Fatalf("Server failed to start: %s", err)
		}
	}
}

func shutdown(server *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Failed to shutdown server gracefully: %s", err)
	}
	log.Println("Server gracefully stopped")
}

//8 must watch
