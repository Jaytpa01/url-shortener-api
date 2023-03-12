package handler

import (
	"net/http"

	"github.com/Jaytpa01/url-shortener-api/api"
	"github.com/Jaytpa01/url-shortener-api/pkg/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

// RedirectToTargetUrl handles redirecting the user
// to the target link from the generated link on our server
func (h *Handler) RedirectToTargetUrl() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := chi.URLParam(r, "token")

		url, err := h.urlService.GetUrlByToken(r.Context(), token)
		if err != nil {
			getUrlErr := api.EnsureApiError(err)

			render.Status(r, getUrlErr.Status())
			render.JSON(w, r, getUrlErr)
			return
		}

		err = h.urlService.IncrementUrlVisits(r.Context(), url)
		if err != nil {
			incrementErr := api.EnsureApiError(err)

			render.Status(r, incrementErr.Status())
			render.JSON(w, r, incrementErr)
			return
		}

		http.Redirect(w, r, url.TargetUrl, http.StatusMovedPermanently)
	}
}

// GetUrlVisits handles fetching the amount of unique vists a generated link has received.
func (h *Handler) GetUrlVisits() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := chi.URLParam(r, "token")

		url, err := h.urlService.GetUrlByToken(r.Context(), token)
		if err != nil {
			getUrlErr := api.EnsureApiError(err)

			render.Status(r, getUrlErr.Status())
			render.JSON(w, r, getUrlErr)
			return
		}

		visitRes := &api.UrlVisitsResponse{
			Visits: url.Visits,
		}

		render.JSON(w, r, visitRes)
	}
}

// ShortenUrl handles returning a shortened url
func (h *Handler) ShortenUrl() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &api.CreateUrlRequest{}

		// decode request payload
		err := h.decoder.DecodeJSON(w, r, req)
		if err != nil {
			decodeErr := api.EnsureApiError(err)

			render.Status(r, decodeErr.Status())
			render.JSON(w, r, decodeErr)
			return
		}

		// create a shortened url
		createdUrl, err := h.urlService.ShortenUrl(r.Context(), req.Url)
		if err != nil {
			createUrlErr := api.EnsureApiError(err)

			render.Status(r, createUrlErr.Status())
			render.JSON(w, r, createUrlErr)
			return
		}

		// convert the data model to an api response model
		apiResponse := &api.UrlResponse{
			Token:     createdUrl.Token,
			TargetUrl: createdUrl.TargetUrl,
			QRCode:    utils.GenerateQRCodeLink(createdUrl.TargetUrl),
		}

		// return the successfully created url with HTTP Status Created
		render.Status(r, http.StatusCreated)
		render.JSON(w, r, apiResponse)
	}
}

// TODO: Write tests for this handler
// LengthenUrl handles returning a longer url, this is a gimmick endpoint
func (h *Handler) LengthenUrl() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &api.CreateUrlRequest{}

		err := h.decoder.DecodeJSON(w, r, req)
		if err != nil {
			decodeErr := api.EnsureApiError(err)

			render.Status(r, decodeErr.Status())
			render.JSON(w, r, decodeErr)
			return
		}

		createdUrl, err := h.urlService.LengthenUrl(r.Context(), req.Url)
		if err != nil {
			createUrlErr := api.EnsureApiError(err)

			render.Status(r, createUrlErr.Status())
			render.JSON(w, r, createUrlErr)
			return
		}

		// convert the data model to an api response model
		apiResponse := &api.UrlResponse{
			Token:     createdUrl.Token,
			TargetUrl: createdUrl.TargetUrl,
			QRCode:    utils.GenerateQRCodeLink(createdUrl.TargetUrl),
		}

		render.Status(r, http.StatusCreated)
		render.JSON(w, r, apiResponse)
	}
}
