package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Jaytpa01/url-shortener-api/api"
	"github.com/Jaytpa01/url-shortener-api/internal/handler"
	"github.com/Jaytpa01/url-shortener-api/internal/repository"
	"github.com/Jaytpa01/url-shortener-api/internal/service"
	"github.com/Jaytpa01/url-shortener-api/pkg/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httprate"
	"github.com/go-chi/render"
)

func serveHTTP() {
	logger := logger.NewApiLogger()

	// create our router and setup rate limiting
	router := chi.NewRouter()
	router.Use(httprate.Limit(1, time.Second, httprate.WithLimitHandler(func(w http.ResponseWriter, r *http.Request) {
		limitErr := api.NewTooManyRequests()
		render.Status(r, http.StatusTooManyRequests)
		render.JSON(w, r, limitErr)
	})))

	// create our repo(s)
	urlRepo := repository.NewInMemoryRepo()

	// create our service(s)
	urlService := service.NewUrlService(&service.Config{
		Logger:  logger,
		UrlRepo: urlRepo,
	})

	// create a new handler for our services.
	// This handles mapping routes to their services
	handler.NewHandler(&handler.Config{
		Router:     router,
		UrlService: urlService,
	})

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
