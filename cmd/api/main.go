package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-api-practice/internal/config"
	"go-api-practice/internal/database"
	"go-api-practice/internal/router"

	"github.com/joho/godotenv"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("load .env: %w", err)
	}

	databaseConfig := config.LoadDatabase()

	connectContext, cancelConnect := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelConnect()

	db, err := database.OpenPostgres(connectContext, databaseConfig.DSN())
	if err != nil {
		return err
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("close database: %v", err)
		}
	}()

	server := &http.Server{
		Addr:              ":8080",
		Handler:           router.New(db),
		ReadHeaderTimeout: 5 * time.Second,
	}

	serverErrors := make(chan error, 1)
	go func() {
		log.Println("server is listening on http://localhost:8080")
		serverErrors <- server.ListenAndServe()
	}()

	stopContext, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	select {
	case err := <-serverErrors:
		if !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("run HTTP server: %w", err)
		}
		return nil
	case <-stopContext.Done():
		log.Println("shutting down server...")
	}

	shutdownContext, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelShutdown()

	if err := server.Shutdown(shutdownContext); err != nil {
		return fmt.Errorf("shutdown HTTP server: %w", err)
	}

	return nil
}
