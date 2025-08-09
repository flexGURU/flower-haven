package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/flexGURU/flower-haven/backend/internal/handlers"
	"github.com/flexGURU/flower-haven/backend/internal/postgres"
	"github.com/flexGURU/flower-haven/backend/pkg"
)

func main() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// load config
	config, err := pkg.LoadConfig("/home/emilio-cliff/flower-haven/backend/.envs/.local")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	tokenMaker := pkg.NewJWTMaker(config.TOKEN_SYMMETRIC_KEY, config.TOKEN_ISSUER)
	if err != nil {
		log.Fatalf("Error creating token maker: %v", err)
	}

	// open database
	store := postgres.NewStore(config)
	err = store.OpenDB(context.Background())
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	// initialize repo
	postgresRepo := postgres.NewPostgresRepo(store)

	// start server
	server := handlers.NewServer(config, tokenMaker, postgresRepo)

	log.Println("starting server at address: ", config.SERVER_ADDRESS)
	if err := server.Start(); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

	// token, _ := tokenMaker.CreateToken(1, "test@test.com", true, 10*time.Hour)
	// log.Println("token: ", token)

	<-quit
	signal.Stop(quit)

	log.Println("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Stop(ctx); err != nil {
		log.Fatalf("Error stopping server: %v", err)
	}

	os.Exit(0)
}
