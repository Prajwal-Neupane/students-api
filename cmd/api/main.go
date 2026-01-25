package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Prajwal-Neupane/students-api/internal/config"
	"github.com/Prajwal-Neupane/students-api/internal/handler/student"
	"github.com/Prajwal-Neupane/students-api/internal/storage/sqlite"
)


func main() {
// load config
	cfg := config.MustLoad()
// database setup

	storage, err:= sqlite.New(cfg)

	if err != nil {
		log.Fatal(err)
	}
	slog.Info("Storage INitialized", slog.String("env", cfg.Env), slog.String("version", "1.0.0"))

// setup router

	router := http.NewServeMux()


	router.HandleFunc("POST /api/students", student.New(storage))
	router.HandleFunc("GET /api/students/{id}", student.GetById(storage))
// setup server

	server := http.Server{
		Addr: cfg.Address,
		Handler: router,
	}

	// fmt.Println("Server started")
	slog.Info("Server started", slog.String("address", cfg.Address))
	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {

	err := server.ListenAndServe()

	if err != nil {
		log.Fatal("failed to start server", err)
	}
	}()

	<-done

	slog.Info("Shutting down the server")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5)
	defer cancel()

	server.Shutdown(ctx)

	// if err != nil {
	// 	slog.Error("Failed to shutdown server", slog.String("error", err.Error()))
	// }

	// slog.Info("Server shutdown successfully");

	
}