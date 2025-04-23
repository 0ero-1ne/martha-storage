package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/0ero-1ne/martha-storage/config"
	"github.com/0ero-1ne/martha-storage/db"
	"github.com/0ero-1ne/martha-storage/router"
	"github.com/0ero-1ne/martha-storage/server"
)

const configPath = "config/ini/config.ini"

func main() {
	config, err := config.Init(configPath)

	if err != nil {
		panic("Config initialization error: " + err.Error())
	}

	database, err := db.NewDbConnection(config.DatabaseConfig)

	if err != nil {
		panic("Can not connect to database: " + err.Error())
	}

	router := router.NewRouter(*config, database)
	server := server.NewHttpServer(config.ServerConfig, router)

	go func() {
		log.Printf("Server is listening on %s\n", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic("Server listening error: " + err.Error())
		}
	}()

	done := make(chan bool, 1)

	go gracefulShutdown(server, done)

	<-done
	log.Println("Graceful shutdown complete.")
}

func gracefulShutdown(server *http.Server, done chan bool) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	log.Println("Shutting down gracefully, press Ctrl+C again to force")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Println("Server exiting")
	done <- true
}
