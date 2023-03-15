package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Jaytpa01/url-shortener-api/config"
	"github.com/Jaytpa01/url-shortener-api/internal/handler"
	"github.com/Jaytpa01/url-shortener-api/internal/repository"
	"github.com/Jaytpa01/url-shortener-api/internal/service"
	"github.com/Jaytpa01/url-shortener-api/pkg/logger"
	"github.com/Jaytpa01/url-shortener-api/pkg/utils"
	"github.com/go-chi/chi/v5"
)

func serveHTTP() {
	logger := logger.NewApiLogger()

	configPath := utils.GetConfigFilepathFromFilename("config.local.yaml")
	config, err := config.LoadConfig(configPath)
	if err != nil {
		logger.Fatalf("couldn't load api configuratiom: %v", err)
		return
	}

	// create our repo(s)
	urlRepo := repository.NewInMemoryRepo()

	// create our service(s)
	urlService := service.NewUrlService(&service.Config{
		Logger:  logger,
		UrlRepo: urlRepo,
	})

	// create our router
	router := chi.NewRouter()

	// create a new handler for our services.
	cfg := &handler.Config{
		Router:     router,
		ApiConfig:  config,
		UrlService: urlService,
	}
	err = handler.NewHandler(cfg)
	if err != nil {
		logger.Fatalf("couldn't create a valid api handler: %v", err)
	}

	httpServer := &http.Server{
		Addr:    ":8080", // TODO: get port from config file
		Handler: router,
	}

	// start the server
	go func() {
		if err := httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			logger.Fatalf("HTTP server error: %v", err)
		}
		logger.Info("Stopped serving new connections.")
	}()

	logger.Info("Server Initialised.")

	// Wait for a termination signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sigChan

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		logger.Fatalf("HTTP server shutdown error: %v", err)
	}
	logger.Info("Graceful shutdown complete.")
}
