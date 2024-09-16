package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725725133-team-77957/zadanie-6105/src/cmd"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725725133-team-77957/zadanie-6105/src/internal/api/handlers"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725725133-team-77957/zadanie-6105/src/internal/database"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725725133-team-77957/zadanie-6105/src/internal/repository"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	//_ "github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen"
)

//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -config ../../api/cfg.yaml ../../api/openapi.yaml

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	config, err := cmd.Load()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Database.Host,
		config.Database.Port,
		config.Database.User,
		config.Database.Password,
		config.Database.Name,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Error opening database connection: %v", err)
	}
	defer func() { _ = db.Close() }()

	database := database.NewDataBase(db)
	resp := repository.NewRepository(database)
	server := handlers.NewServer(resp)
	mux := http.NewServeMux()
	h := handlers.HandlerFromMuxWithBaseURL(server, mux, "/api")

	s := &http.Server{
		Addr:    config.ServerAddress,
		Handler: h,
	}

	// Start server in a new goroutine
	go func() {
		log.Printf("Service started on port %s", config.ServerAddress)
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Error while starting server: %v", err)
		}
	}()

	// Wait for termination signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit

	// Shutdown server gracefully
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		log.Printf("Error during server shutdown: %v", err)
	}
	log.Println("Server gracefully stopped")
}
