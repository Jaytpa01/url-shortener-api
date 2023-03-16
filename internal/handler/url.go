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
func (h *handler) RedirectToTargetUrl() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := chi.URLParam(r, "token")

		url, err := h.urlService.FindUrlByToken(r.Context(), token)
		if err != nil {
			returnApiError(w, r, err)
			return
		}

		err = h.urlService.IncrementUrlVisits(r.Context(), url)
		if err != nil {
			returnApiError(w, r, err)
			return
		}

		http.Redirect(w, r, url.TargetUrl, http.StatusMovedPermanently)
	}
}

// GetUrlVisits handles fetching the amount of unique vists a generated link has received.
func (h *handler) GetUrlVisits() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := chi.URLParam(r, "token")

		url, err := h.urlService.FindUrlByToken(r.Context(), token)
		if err != nil {
			returnApiError(w, r, err)
			return
		}

		visitRes := &api.UrlVisitsResponse{
			Visits: url.Visits,
		}

		render.JSON(w, r, visitRes)
	}
}

// ShortenUrl handles returning a shortened url
func (h *handler) ShortenUrl() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &api.CreateUrlRequest{}

		// decode request payload
		err := h.decoder.DecodeJSON(w, r, req)
		if err != nil {
			returnApiError(w, r, err)
			return
		}

		// create a shortened url
		createdUrl, err := h.urlService.ShortenUrl(r.Context(), req.Url)
		if err != nil {
			returnApiError(w, r, err)
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
func (h *handler) LengthenUrl() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &api.CreateUrlRequest{}

		err := h.decoder.DecodeJSON(w, r, req)
		if err != nil {
			returnApiError(w, r, err)
			return
		}

		createdUrl, err := h.urlService.LengthenUrl(r.Context(), req.Url)
		if err != nil {
			returnApiError(w, r, err)
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

// GetAllUrls is a handler only available in a development environment. It
// gets all urls in the repo and returns them
func (h *handler) GetAllUrls() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urls, err := h.urlService.GetAllUrls(r.Context())
		if err != nil {
			returnApiError(w, r, err)
			return
		}

		render.JSON(w, r, urls)
	}
}

// returnApiError is a helper function to ensure any errors we return to the client conform to our standard error response
func returnApiError(w http.ResponseWriter, r *http.Request, err error) {
	apiErr := api.EnsureApiError(err)
	render.Status(r, apiErr.Status())
	render.JSON(w, r, apiErr)
}
