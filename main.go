package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	log.Println("Hello, it's a me, twotter!")

	ctx := context.Background()
	if err := run(ctx); err != nil {
		log.Fatal("twotter terminated abnormally: ", err)
	}

	log.Println("Ok lady, I love you buhbye!")
}

func run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	model := LoadTwots()

	server, err := NewServer(ctx, model)
	if err != nil {
		return fmt.Errorf("Failed to create server: %w", err)
	}
	go shutdownHandler(ctx, server)

	log.Println("Listening on :8080")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("Failed to listen and serve: %w", err)
	}

	return nil
}

func shutdownHandler(ctx context.Context, server *http.Server) error {
	<-ctx.Done()

	log.Println("Shutting down...")
	shutdownCtx := context.Background()
	shutdownCtx, cancel := context.WithTimeout(shutdownCtx, 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatal("Failed to shutdown: ", err)
	}

	return nil
}
