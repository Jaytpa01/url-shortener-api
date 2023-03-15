package handler

import (
	"net/http"
	"time"

	"github.com/Jaytpa01/url-shortener-api/api"
	"github.com/Jaytpa01/url-shortener-api/config"
	"github.com/Jaytpa01/url-shortener-api/internal/service"
	"github.com/Jaytpa01/url-shortener-api/pkg/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"github.com/go-chi/render"
)

// Config is used to setup the router and services, which
// are then injected into the handler.
type Config struct {
	Router     *chi.Mux
	Decoder    utils.JSONDecoder
	ApiConfig  *config.Config
	UrlService service.UrlService
}

// validate is used to validate the handler config based off of the current environment (i.e development, production)
// TODO: implement validating config
func (c *Config) validate() error {
	return nil
}

// Handler holds the required services for the handler and app to function
type handler struct {
	router     *chi.Mux
	decoder    utils.JSONDecoder
	apiConfig  *config.Config
	urlService service.UrlService
}

// NewHandler initialises the handler with the injected services, and sets up the http routes.
// It doesn't return a handler instance, only an error, as it directly deals with the chi router.
// If there is an error with the configuration, an error is returned.
func NewHandler(cfg *Config) error {
	if err := cfg.validate(); err != nil {
		return err
	}

	// if we haven't injected a decoder, use our implementation
	decoder := cfg.Decoder
	if decoder == nil {
		decoder = utils.NewJSONDecoder()
	}

	// create a handler
	h := newHandler(cfg.Router, decoder, cfg.ApiConfig, cfg.UrlService)

	// get a reference to the router and
	// put it in a variable easier to work with
	r := h.router

	r.Use(middleware.Logger)

	r.Use(httprate.Limit(10, time.Second, httprate.WithLimitHandler(func(w http.ResponseWriter, r *http.Request) {
		limitErr := api.NewTooManyRequests()
		render.Status(r, http.StatusTooManyRequests)
		render.JSON(w, r, limitErr)
	})))

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		MaxAge:         300, // Maximum value not ignored by any of major browsers
	}))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, map[string]string{
			"status": "ok",
		})
	})
	r.Get("/{token}", h.RedirectToTargetUrl())
	r.Get("/{token}/visits", h.GetUrlVisits())
	r.Post("/shorten", h.ShortenUrl())
	r.Post("/lengthen", h.LengthenUrl())

	if h.apiConfig != nil {
		switch h.apiConfig.Server.Environment {
		case "dev", "development", "local":
			r.Get("/all", h.GetAllUrls())
		}
	}

	return nil
}

// new handler is a package scoped facotry function for creating a handler
func newHandler(router *chi.Mux, decoder utils.JSONDecoder, apiConfig *config.Config, urlService service.UrlService) *handler {
	return &handler{
		router:     router,
		decoder:    decoder,
		apiConfig:  apiConfig,
		urlService: urlService,
	}
}
