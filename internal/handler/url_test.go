package handler

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Jaytpa01/url-shortener-api/api"
	"github.com/Jaytpa01/url-shortener-api/internal/entity"
	"github.com/Jaytpa01/url-shortener-api/internal/mocks"
	"github.com/Jaytpa01/url-shortener-api/pkg/utils"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	mockEmptyContext  = "*context.emptyCtx"
	contentTypeHeader = "Content-Type"
	contentTypeJSON   = "application/json"
	exampleUrl        = "https://example.com"
)

func TestHandler_Url_RedirectToTargetUrl(t *testing.T) {
	// setup
	token := utils.GenerateRandomString(6)
	req := httptest.NewRequest(http.MethodGet, "/"+token, nil)
	rec := httptest.NewRecorder()

	mockCreatedAt := time.Now()
	mockResponse := &entity.Url{
		Token:     token,
		TargetUrl: exampleUrl,
		Visits:    0,
		CreatedAt: mockCreatedAt,
	}

	mockUrlService := mocks.NewMockUrlService()
	mockUrlService.On("GetUrlByToken", mock.Anything, token).Return(mockResponse, nil)
	mockUrlService.On("IncrementUrlVisits", mock.Anything, mock.Anything).Return(nil)
	mockResponse.Visits++

	r := chi.NewRouter()
	NewHandler(&Config{
		Router:     r,
		UrlService: mockUrlService,
	})

	r.ServeHTTP(rec, req)
	result := rec.Result()

	redirectLocation, _ := result.Location()

	// Assertions
	assert.Equal(t, http.StatusMovedPermanently, result.StatusCode)
	assert.Equal(t, exampleUrl, redirectLocation.String())
}

func TestHandler_Url_RedirectToTargetUrl_UrlDoesntExist(t *testing.T) {
	// setup
	token := utils.GenerateRandomString(6)
	req := httptest.NewRequest(http.MethodGet, "/"+token, nil)
	rec := httptest.NewRecorder()

	mockMsg := fmt.Sprintf("Couldn't find URL with token (%s).", token)
	mockResponse := api.NewNotFound("url-not-found", mockMsg)

	mockUrlService := mocks.NewMockUrlService()
	mockUrlService.On("GetUrlByToken", mock.Anything, token).Return(nil, mockResponse)

	r := chi.NewRouter()
	NewHandler(&Config{
		Router:     r,
		UrlService: mockUrlService,
	})

	r.ServeHTTP(rec, req)
	result := rec.Result()

	// Assertions
	assert.Equal(t, http.StatusNotFound, result.StatusCode)
	assert.Equal(t, fmt.Sprintf("{\"type\":\"NOT_FOUND\",\"code\":\"url-not-found\",\"message\":\"%s\"}", mockMsg), strings.Trim(rec.Body.String(), "\n"))
}

func TestHandler_Url_RedirectToTargetUrl_UnknownErrorOccured(t *testing.T) {
	// setup
	token := utils.GenerateRandomString(6)
	req := httptest.NewRequest(http.MethodGet, "/"+token, nil)
	rec := httptest.NewRecorder()

	unknownErr := errors.New("some unknown error")

	mockUrlService := mocks.NewMockUrlService()
	mockUrlService.On("GetUrlByToken", mock.Anything, token).Return(nil, unknownErr)

	r := chi.NewRouter()
	NewHandler(&Config{
		Router:     r,
		UrlService: mockUrlService,
	})

	r.ServeHTTP(rec, req)
	result := rec.Result()

	// Assertions
	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
	assert.Equal(t, "{\"type\":\"INTERNAL\",\"code\":\"unknown\",\"message\":\"An internal server error occured.\"}", strings.Trim(rec.Body.String(), "\n"))

	// Ensure apiError.Debug is correct
	apiError := api.EnsureApiError(unknownErr)
	assert.Equal(t, unknownErr.Error(), apiError.Debug)
}

func TestHandler_Url_RedirectToTargetUrl_FailedToIncrementVisits(t *testing.T) {
	// setup
	token := utils.GenerateRandomString(6)
	req := httptest.NewRequest(http.MethodGet, "/"+token, nil)
	rec := httptest.NewRecorder()

	mockCreatedAt := time.Now()
	mockResponse := &entity.Url{
		Token:     token,
		TargetUrl: exampleUrl,
		Visits:    0,
		CreatedAt: mockCreatedAt,
	}

	unknownErr := errors.New("some unknown error")

	mockUrlService := mocks.NewMockUrlService()
	mockUrlService.On("GetUrlByToken", mock.Anything, token).Return(mockResponse, nil)
	mockUrlService.On("IncrementUrlVisits", mock.Anything, mock.Anything).Return(unknownErr)

	r := chi.NewRouter()
	NewHandler(&Config{
		Router:     r,
		UrlService: mockUrlService,
	})

	r.ServeHTTP(rec, req)
	result := rec.Result()

	// Assertions
	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
	assert.Equal(t, "{\"type\":\"INTERNAL\",\"code\":\"unknown\",\"message\":\"An internal server error occured.\"}", strings.Trim(rec.Body.String(), "\n"))

	// Ensure apiError.Debug is correct
	apiError := api.EnsureApiError(unknownErr)
	assert.Equal(t, unknownErr.Error(), apiError.Debug)
}

func TestHandler_Url_GetUrlVisits(t *testing.T) {
	// setup
	token := utils.GenerateRandomString(6)
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s/visits", token), nil)
	rec := httptest.NewRecorder()

	mockCreatedAt := time.Now()
	urlVisits := 42
	mockResponse := &entity.Url{
		Token:     token,
		TargetUrl: exampleUrl,
		Visits:    urlVisits,
		CreatedAt: mockCreatedAt,
	}

	mockUrlService := mocks.NewMockUrlService()
	mockUrlService.On("GetUrlByToken", mock.Anything, token).Return(mockResponse, nil)

	r := chi.NewRouter()
	NewHandler(&Config{
		Router:     r,
		UrlService: mockUrlService,
	})

	r.ServeHTTP(rec, req)
	result := rec.Result()

	// Assertions
	assert.Equal(t, http.StatusOK, result.StatusCode)
	assert.Equal(t, fmt.Sprintf("{\"visits\":%d}", urlVisits), strings.Trim(rec.Body.String(), "\n"))
}

func TestHandler_Url_GetUrlVisits_UrlDoesntExist(t *testing.T) {
	// setup
	token := utils.GenerateRandomString(6)
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s/visits", token), nil)
	rec := httptest.NewRecorder()

	mockMsg := fmt.Sprintf("Couldn't find URL with token (%s).", token)
	mockResponse := api.NewNotFound("url-not-found", mockMsg)

	mockUrlService := mocks.NewMockUrlService()
	mockUrlService.On("GetUrlByToken", mock.Anything, token).Return(nil, mockResponse)

	r := chi.NewRouter()
	NewHandler(&Config{
		Router:     r,
		UrlService: mockUrlService,
	})

	r.ServeHTTP(rec, req)
	result := rec.Result()

	// Assertions
	assert.Equal(t, http.StatusNotFound, result.StatusCode)
	assert.Equal(t, fmt.Sprintf("{\"type\":\"NOT_FOUND\",\"code\":\"url-not-found\",\"message\":\"%s\"}", mockMsg), strings.Trim(rec.Body.String(), "\n"))
}

func TestHandler_Url_ShortenUrl(t *testing.T) {
	// setup
	reqBody := fmt.Sprintf("{\"url\":\"%s\"}", exampleUrl)
	req := httptest.NewRequest(http.MethodPost, "/shorten", strings.NewReader(reqBody))
	req.Header.Set(contentTypeHeader, contentTypeJSON)
	rec := httptest.NewRecorder()

	token := utils.GenerateRandomString(6)
	mockUrlResponse := &entity.Url{
		Token:     token,
		TargetUrl: exampleUrl,
		Visits:    0,
		CreatedAt: time.Now(),
	}

	mockUrlService := mocks.NewMockUrlService()
	mockUrlService.On("ShortenUrl", mock.Anything, exampleUrl).Return(mockUrlResponse, nil)

	r := chi.NewRouter()
	NewHandler(&Config{
		Router:     r,
		UrlService: mockUrlService,
	})

	r.ServeHTTP(rec, req)
	result := rec.Result()

	// Assertions
	assert.Equal(t, http.StatusCreated, result.StatusCode)
	assert.Equal(t, fmt.Sprintf("{\"token\":\"%s\",\"target_url\":\"%s\",\"qr_code\":\"%s\"}", token, exampleUrl, utils.GenerateQRCodeLink(exampleUrl)), strings.Trim(rec.Body.String(), "\n"))
}

func TestHandler_Url_ShortenUrl_UnknownFieldError(t *testing.T) {
	reqBody := fmt.Sprintf("{\"url\":\"%s\",\"bad_field\":\"some dummy data\"}", exampleUrl)
	req := httptest.NewRequest(http.MethodPost, "/shorten", strings.NewReader(reqBody))
	req.Header.Set(contentTypeHeader, contentTypeJSON)
	rec := httptest.NewRecorder()

	mockUrlService := mocks.NewMockUrlService()

	r := chi.NewRouter()
	NewHandler(&Config{
		Router:     r,
		UrlService: mockUrlService,
	})

	r.ServeHTTP(rec, req)
	result := rec.Result()

	// Assertions
	assert.Equal(t, http.StatusBadRequest, result.StatusCode)
	assert.Equal(t, "{\"type\":\"BAD_REQUEST\",\"code\":\"unknown-field\",\"message\":\"Request body contains unknown field (\\\"bad_field\\\").\"}", strings.Trim(rec.Body.String(), "\n"))
}

func TestHandler_Url_ShortenUrl_EofJSON(t *testing.T) {
	// setup
	reqBody := "{\"url\":"
	req := httptest.NewRequest(http.MethodPost, "/shorten", strings.NewReader(reqBody))
	req.Header.Set(contentTypeHeader, contentTypeJSON)
	rec := httptest.NewRecorder()

	r := chi.NewRouter()
	NewHandler(&Config{
		Router: r,
	})

	r.ServeHTTP(rec, req)
	result := rec.Result()

	// Assertions
	assert.Equal(t, http.StatusBadRequest, result.StatusCode)
	assert.Equal(t, "{\"type\":\"BAD_REQUEST\",\"code\":\"json-eof-error\",\"message\":\"Request body contains badly-formed JSON.\"}", strings.Trim(rec.Body.String(), "\n"))
}

func TestHandler_Url_ShortenUrl_JSONSyntaxError(t *testing.T) {
	// setup
	reqBody := "{\"url\":}"
	req := httptest.NewRequest(http.MethodPost, "/shorten", strings.NewReader(reqBody))
	req.Header.Set(contentTypeHeader, contentTypeJSON)
	rec := httptest.NewRecorder()

	r := chi.NewRouter()
	NewHandler(&Config{
		Router: r,
	})

	r.ServeHTTP(rec, req)
	result := rec.Result()

	// Assertions
	assert.Equal(t, http.StatusBadRequest, result.StatusCode)
	assert.Equal(t, "{\"type\":\"BAD_REQUEST\",\"code\":\"json-syntax-error\",\"message\":\"Request body contains badly-formed JSON (at position 8).\"}", strings.Trim(rec.Body.String(), "\n"))
}
func TestHandler_Url_ShortenUrl_UnknownDecodeJSONError(t *testing.T) {
	// setup
	reqBody := fmt.Sprintf("{\"url\":\"%s\"}", exampleUrl)
	req := httptest.NewRequest(http.MethodPost, "/shorten", strings.NewReader(reqBody))
	req.Header.Set(contentTypeHeader, contentTypeJSON)
	rec := httptest.NewRecorder()

	errUnknownDecodeJSON := errors.New("unknown error decoding json")
	mockJSONDecoder := mocks.NewMockJSONDecoder()
	mockJSONDecoder.On("DecodeJSON", mock.Anything, mock.Anything, mock.Anything).Return(errUnknownDecodeJSON)

	r := chi.NewRouter()
	NewHandler(&Config{
		Router:  r,
		Decoder: mockJSONDecoder,
	})

	r.ServeHTTP(rec, req)
	result := rec.Result()

	// Assertions
	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
	assert.Equal(t, "{\"type\":\"INTERNAL\",\"code\":\"unknown\",\"message\":\"An internal server error occured.\"}", strings.Trim(rec.Body.String(), "\n"))

	// Ensure apiError.Debug is correct
	apiError := api.EnsureApiError(errUnknownDecodeJSON)
	assert.Equal(t, errUnknownDecodeJSON.Error(), apiError.Debug)
}

func TestHandler_Url_ShortenUrl_EmptyRequestBody(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/shorten", nil)
	req.Header.Set(contentTypeHeader, contentTypeJSON)
	rec := httptest.NewRecorder()

	r := chi.NewRouter()
	NewHandler(&Config{
		Router: r,
	})

	r.ServeHTTP(rec, req)
	result := rec.Result()

	// Assertions
	assert.Equal(t, http.StatusBadRequest, result.StatusCode)
	assert.Equal(t, "{\"type\":\"BAD_REQUEST\",\"code\":\"no-empty-requests\",\"message\":\"Request body must not be empty.\"}", strings.Trim(rec.Body.String(), "\n"))
}

func TestHandler_Url_ShortenUrl_FailedToCreate(t *testing.T) {
	// setup
	reqBody := fmt.Sprintf("{\"url\":\"%s\"}", exampleUrl)
	req := httptest.NewRequest(http.MethodPost, "/shorten", strings.NewReader(reqBody))
	req.Header.Set(contentTypeHeader, contentTypeJSON)
	rec := httptest.NewRecorder()

	mockError := errors.New("some error shortening url")

	mockUrlService := mocks.NewMockUrlService()
	mockUrlService.On("ShortenUrl", mock.Anything, exampleUrl).Return(nil, mockError)

	r := chi.NewRouter()
	NewHandler(&Config{
		Router:     r,
		UrlService: mockUrlService,
	})

	r.ServeHTTP(rec, req)
	result := rec.Result()

	// Assertions
	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
	assert.Equal(t, "{\"type\":\"INTERNAL\",\"code\":\"unknown\",\"message\":\"An internal server error occured.\"}", strings.Trim(rec.Body.String(), "\n"))

	// Ensure apiError.Debug is correct
	apiError := api.EnsureApiError(mockError)
	assert.Equal(t, mockError.Error(), apiError.Debug)
}
