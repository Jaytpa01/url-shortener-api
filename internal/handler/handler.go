package handler

import (
	"github.com/Jaytpa01/url-shortener-api/internal/service"
	"github.com/Jaytpa01/url-shortener-api/pkg/utils"
	"github.com/go-chi/chi/v5"
)

// Handler holds the required services for the handler and app to function
type Handler struct {
	decoder    utils.JSONDecoder
	urlService service.UrlService
}

// Config is used to setup the router and services, which
// are then injected into the handler.
type Config struct {
	Router     *chi.Mux
	Decoder    utils.JSONDecoder
	UrlService service.UrlService
}

// NewHandler initialises the handler with the injected services, and
// sets up the http routes.
// It doesn't return anything, as it directly deals with the chi router.
func NewHandler(cfg *Config) {

	// if we haven't injected a decoder, use our implementation
	decoder := cfg.Decoder
	if decoder == nil {
		decoder = utils.NewJSONDecoder()
	}

	handler := &Handler{
		decoder:    decoder,
		urlService: cfg.UrlService,
	}

	// get a reference to the router and
	// put it in a variable easier to work with
	r := cfg.Router

	r.Get("/{token}", handler.RedirectToTargetUrl())
	r.Get("/{token}/visits", handler.GetUrlVisits())
	r.Post("/shorten", handler.ShortenUrl())
	r.Post("/lengthen", handler.LengthenUrl())

}
