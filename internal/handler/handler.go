package handler

import (
	"fmt"

	"github.com/Jaytpa01/url-shortener-api/internal/service"
	"github.com/go-chi/chi/v5"
)

// Handler holds the required services for the handler and app to function
type Handler struct {
	UrlService service.UrlService
}

// Config is used to setup the router and services, which
// are then injected into the handler.
type Config struct {
	Router     *chi.Mux
	UrlService service.UrlService
}

// NewHandler initialises the handler with the injected services, and
// sets up the http routes.
// It doesn't return anything, as it directly deals with the chi router.
func NewHandler(cfg *Config) {
	handler := &Handler{
		UrlService: cfg.UrlService,
	}

	fmt.Println(handler)

	// TODO: implement routes here
}
