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
	"github.com/Jaytpa01/url-shortener-api/config"
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

var (
	apiConfig = &config.Config{Server: config.ServerConfig{Environment: "test"}}
)

func TestHandler_Url_RedirectToTargetUrl(t *testing.T) {
	// setup
	token := utils.NewRandomiser().GenerateRandomString(6)
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
	mockUrlService.On("FindUrlByToken", mock.Anything, token).Return(mockResponse, nil)
	mockUrlService.On("IncrementUrlVisits", mock.Anything, mock.Anything).Return(nil)
	mockResponse.Visits++

	r := chi.NewRouter()
	NewHandler(&Config{
		Router:     r,
		ApiConfig:  apiConfig,
		UrlService: mockUrlService,
	})

	r.ServeHTTP(rec, req)
	result := rec.Result()

	redirectLocation, _ := result.Location()

	// Assertions
	assert.Equal(t, http.StatusMovedPermanently, result.StatusCode)
	assert.Equal(t, exampleUrl, redirectLocation.String())
}

func TestHandler_Url_GetIndex(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	r := chi.NewRouter()
	NewHandler(&Config{
		Router:    r,
		ApiConfig: apiConfig,
	})

	r.ServeHTTP(rec, req)

	assert.Equal(t, "{\"status\":\"ok\"}", strings.Trim(rec.Body.String(), "\n"))
}

func TestHandler_Url_RedirectToTargetUrl_UrlDoesntExist(t *testing.T) {
	// setup
	token := utils.NewRandomiser().GenerateRandomString(6)
	req := httptest.NewRequest(http.MethodGet, "/"+token, nil)
	rec := httptest.NewRecorder()

	mockMsg := fmt.Sprintf("Couldn't find URL with token (%s).", token)
	mockResponse := api.NewNotFound("url-not-found", mockMsg)

	mockUrlService := mocks.NewMockUrlService()
	mockUrlService.On("FindUrlByToken", mock.Anything, token).Return(nil, mockResponse)

	r := chi.NewRouter()
	NewHandler(&Config{
		Router:     r,
		ApiConfig:  apiConfig,
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
	token := utils.NewRandomiser().GenerateRandomString(6)
	req := httptest.NewRequest(http.MethodGet, "/"+token, nil)
	rec := httptest.NewRecorder()

	unknownErr := errors.New("some unknown error")

	mockUrlService := mocks.NewMockUrlService()
	mockUrlService.On("FindUrlByToken", mock.Anything, token).Return(nil, unknownErr)

	r := chi.NewRouter()
	NewHandler(&Config{
		Router:     r,
		ApiConfig:  apiConfig,
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
	token := utils.NewRandomiser().GenerateRandomString(6)
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
	mockUrlService.On("FindUrlByToken", mock.Anything, token).Return(mockResponse, nil)
	mockUrlService.On("IncrementUrlVisits", mock.Anything, mock.Anything).Return(unknownErr)

	r := chi.NewRouter()
	NewHandler(&Config{
		Router:     r,
		ApiConfig:  apiConfig,
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
	token := utils.NewRandomiser().GenerateRandomString(6)
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
	mockUrlService.On("FindUrlByToken", mock.Anything, token).Return(mockResponse, nil)

	r := chi.NewRouter()
	NewHandler(&Config{
		Router:     r,
		ApiConfig:  apiConfig,
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
	token := utils.NewRandomiser().GenerateRandomString(6)
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s/visits", token), nil)
	rec := httptest.NewRecorder()

	mockMsg := fmt.Sprintf("Couldn't find URL with token (%s).", token)
	mockResponse := api.NewNotFound("url-not-found", mockMsg)

	mockUrlService := mocks.NewMockUrlService()
	mockUrlService.On("FindUrlByToken", mock.Anything, token).Return(nil, mockResponse)

	r := chi.NewRouter()
	NewHandler(&Config{
		Router:     r,
		ApiConfig:  apiConfig,
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

	token := utils.NewRandomiser().GenerateRandomString(6)
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
		ApiConfig:  apiConfig,
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
		ApiConfig:  apiConfig,
		UrlService: mockUrlService,
	})

	r.ServeHTTP(rec, req)
	result := rec.Result()

	// Assertions
	assert.Equal(t, http.StatusBadRequest, result.StatusCode)
	assert.Equal(t, "{\"type\":\"BAD_REQUEST\",\"code\":\"json/unknown-field\",\"message\":\"Request body contains unknown field (\\\"bad_field\\\").\"}", strings.Trim(rec.Body.String(), "\n"))
}

func TestHandler_Url_ShortenUrl_EofJSON(t *testing.T) {
	// setup
	reqBody := "{\"url\":"
	req := httptest.NewRequest(http.MethodPost, "/shorten", strings.NewReader(reqBody))
	req.Header.Set(contentTypeHeader, contentTypeJSON)
	rec := httptest.NewRecorder()

	r := chi.NewRouter()
	NewHandler(&Config{
		Router:    r,
		ApiConfig: apiConfig,
	})

	r.ServeHTTP(rec, req)
	result := rec.Result()

	// Assertions
	assert.Equal(t, http.StatusBadRequest, result.StatusCode)
	assert.Equal(t, "{\"type\":\"BAD_REQUEST\",\"code\":\"json/eof-error\",\"message\":\"Request body contains badly-formed JSON.\"}", strings.Trim(rec.Body.String(), "\n"))
}

func TestHandler_Url_ShortenUrl_JSONSyntaxError(t *testing.T) {
	// setup
	reqBody := "{\"url\":}"
	req := httptest.NewRequest(http.MethodPost, "/shorten", strings.NewReader(reqBody))
	req.Header.Set(contentTypeHeader, contentTypeJSON)
	rec := httptest.NewRecorder()

	r := chi.NewRouter()
	NewHandler(&Config{
		Router:    r,
		ApiConfig: apiConfig,
	})

	r.ServeHTTP(rec, req)
	result := rec.Result()

	// Assertions
	assert.Equal(t, http.StatusBadRequest, result.StatusCode)
	assert.Equal(t, "{\"type\":\"BAD_REQUEST\",\"code\":\"json/syntax-error\",\"message\":\"Request body contains badly-formed JSON (at position 8).\"}", strings.Trim(rec.Body.String(), "\n"))
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
		Router:    r,
		ApiConfig: apiConfig,
		Decoder:   mockJSONDecoder,
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
		Router:    r,
		ApiConfig: apiConfig,
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
		ApiConfig:  apiConfig,
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

func TestHandler_Url_LengthenUrl(t *testing.T) {
	testCases := []struct {
		name                       string
		requestBody                string
		url                        string
		supplyInvalidContentHeader bool

		useMockDecoder bool
		decoderErr     error

		serviceResponseUrl *entity.Url
		serviceErr         error

		expectedResponseStatus int
		expectedResponseBody   string
	}{
		{
			name:                   "Success",
			requestBody:            "{\"url\":\"https://example.com\"}",
			url:                    "https://example.com",
			serviceResponseUrl:     &entity.Url{Token: "987654", TargetUrl: "https://example.com"},
			expectedResponseStatus: http.StatusCreated,
			expectedResponseBody:   fmt.Sprintf("{\"token\":\"987654\",\"target_url\":\"https://example.com\",\"qr_code\":\"%s\"}", utils.GenerateQRCodeLink("https://example.com")),
		},
		{
			name:                   "Known Service Error",
			requestBody:            "{\"url\":\"https://example.com\"}",
			url:                    "https://example.com",
			serviceErr:             api.NewBadRequest("oops/hahah", "haha oops"),
			expectedResponseStatus: http.StatusBadRequest,
			expectedResponseBody:   "{\"type\":\"BAD_REQUEST\",\"code\":\"oops/hahah\",\"message\":\"haha oops\"}",
		},
		{
			name:                   "Unknown Service Error",
			requestBody:            "{\"url\":\"https://example.com\"}",
			url:                    "https://example.com",
			serviceErr:             errors.New("some unknown error"),
			expectedResponseStatus: http.StatusInternalServerError,
			expectedResponseBody:   "{\"type\":\"INTERNAL\",\"code\":\"unknown\",\"message\":\"An internal server error occured.\"}",
		},
		{
			name:                       "Invalid Content-Type Header",
			requestBody:                "hello",
			supplyInvalidContentHeader: true,
			expectedResponseStatus:     http.StatusUnsupportedMediaType,
		},
		{
			name:        "Malformed JSON",
			requestBody: "{\"url\"}",

			expectedResponseStatus: http.StatusBadRequest,
			expectedResponseBody:   "{\"type\":\"BAD_REQUEST\",\"code\":\"json/syntax-error\",\"message\":\"Request body contains badly-formed JSON (at position 7).\"}",
		},
		{
			name:        "Malformed JSON EOF",
			requestBody: "{\"ur",

			expectedResponseStatus: http.StatusBadRequest,
			expectedResponseBody:   "{\"type\":\"BAD_REQUEST\",\"code\":\"json/eof-error\",\"message\":\"Request body contains badly-formed JSON.\"}",
		},
		{
			name:                   "Unknown Failed To Decode JSON Error",
			requestBody:            "{\"url\":\"https://example.com\"}",
			useMockDecoder:         true,
			decoderErr:             errors.New("some decode error"),
			expectedResponseStatus: http.StatusInternalServerError,
			expectedResponseBody:   "{\"type\":\"INTERNAL\",\"code\":\"unknown\",\"message\":\"An internal server error occured.\"}",
		},
		{
			name:                   "Invalid JSON Field",
			requestBody:            "{\"url\":\"https://example.com\",\"some_other_field\":\"oh boy i sure am some random field\"}",
			expectedResponseStatus: http.StatusBadRequest,
			expectedResponseBody:   "{\"type\":\"BAD_REQUEST\",\"code\":\"json/unknown-field\",\"message\":\"Request body contains unknown field (\\\"some_other_field\\\").\"}",
		},
		{
			name:                   "Empty Request Body",
			expectedResponseStatus: http.StatusBadRequest,
			expectedResponseBody:   "{\"type\":\"BAD_REQUEST\",\"code\":\"no-empty-requests\",\"message\":\"Request body must not be empty.\"}",
		},
		{
			name:                   "Invalid JSON Value for Field",
			requestBody:            "{\"url\":2}",
			expectedResponseStatus: http.StatusBadRequest,
			expectedResponseBody:   "{\"type\":\"BAD_REQUEST\",\"code\":\"json/invalid-field-value\",\"message\":\"Request body contains an invalid value for the \\\"url\\\" field (at position 8).\"}",
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// setup
			var req *http.Request
			if test.requestBody != "" {
				reader := strings.NewReader(test.requestBody)
				req = httptest.NewRequest(http.MethodPost, "/lengthen", reader)
			} else {
				req = httptest.NewRequest(http.MethodPost, "/lengthen", nil)
			}

			if !test.supplyInvalidContentHeader {
				req.Header.Set(contentTypeHeader, contentTypeJSON)
			}
			rec := httptest.NewRecorder()

			mockUrlService := mocks.NewMockUrlService()
			mockUrlService.On("LengthenUrl", mock.Anything, test.url).Return(test.serviceResponseUrl, test.serviceErr)

			r := chi.NewRouter()
			handlerConfig := &Config{Router: r, UrlService: mockUrlService, ApiConfig: apiConfig}
			if test.useMockDecoder {
				mockDecoder := mocks.NewMockJSONDecoder()
				mockDecoder.On("DecodeJSON", mock.Anything, mock.Anything, mock.AnythingOfType("*api.CreateUrlRequest")).Return(test.decoderErr)
				handlerConfig.Decoder = mockDecoder
			} else {
				handlerConfig.Decoder = utils.NewJSONDecoder()
			}

			NewHandler(handlerConfig)
			r.ServeHTTP(rec, req)
			result := rec.Result()

			// Assertions
			assert.Equal(t, test.expectedResponseStatus, result.StatusCode)
			assert.Equal(t, test.expectedResponseBody, strings.Trim(rec.Body.String(), "\n"))
		})
	}
}

func TestHandler_Url_GetAllUrls(t *testing.T) {
	testCases := []struct {
		name string

		serviceResp []entity.Url
		serviceErr  error

		expectedResponseStatus int
		expectedResponseBody   string
	}{
		{
			name:                   "Success",
			serviceResp:            []entity.Url{},
			expectedResponseStatus: http.StatusOK,
			expectedResponseBody:   "[]",
		},
		{
			name:                   "Known Service Error",
			serviceErr:             api.NewBadRequest("bad", "bad message"),
			expectedResponseStatus: http.StatusBadRequest,
			expectedResponseBody:   "{\"type\":\"BAD_REQUEST\",\"code\":\"bad\",\"message\":\"bad message\"}",
		},
		{
			name:                   "Unknown Service Error",
			serviceErr:             errors.New("whoops"),
			expectedResponseStatus: http.StatusInternalServerError,
			expectedResponseBody:   "{\"type\":\"INTERNAL\",\"code\":\"unknown\",\"message\":\"An internal server error occured.\"}",
		},
	}

	for _, test := range testCases {
		// setup
		req := httptest.NewRequest(http.MethodGet, "/all", nil)
		rec := httptest.NewRecorder()

		mockUrlService := mocks.NewMockUrlService()
		mockUrlService.On("GetAllUrls", mock.Anything).Return(test.serviceResp, test.serviceErr)

		r := chi.NewRouter()
		handlerConfig := &Config{Router: r, UrlService: mockUrlService, ApiConfig: &config.Config{Server: config.ServerConfig{Environment: "dev"}}}
		NewHandler(handlerConfig)
		r.ServeHTTP(rec, req)
		result := rec.Result()

		// Assertions
		assert.Equal(t, test.expectedResponseStatus, result.StatusCode)
		assert.Equal(t, test.expectedResponseBody, strings.Trim(rec.Body.String(), "\n"))
	}
}

func TestHandler_NewHandler(t *testing.T) {
	testCases := []struct {
		name          string
		config        *Config
		expectedError error
	}{
		{
			name: "Valid Config",
			config: &Config{
				Router:     chi.NewRouter(),
				ApiConfig:  &config.Config{},
				UrlService: mocks.NewMockUrlService(),
			},
			expectedError: nil,
		},
		{
			name:          "Invalid ApiConfig",
			config:        &Config{},
			expectedError: errors.New("ApiConfig was nil"),
		},
		{
			name: "Invalid Router",
			config: &Config{
				ApiConfig: &config.Config{},
			},
			expectedError: errors.New("router was nil"),
		},
		{
			name: "Invalid Service",
			config: &Config{
				Router:    chi.NewRouter(),
				ApiConfig: &config.Config{},
			},
			expectedError: errors.New("UrlService was nil"),
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			err := NewHandler(test.config)
			assert.Equal(t, test.expectedError, err)
		})
	}
}

func TestHandler_RateLimit(t *testing.T) {
	testCases := []struct {
		name      string
		respCodes []int
	}{
		{
			name:      "Not Reach Limit",
			respCodes: []int{200, 200, 200},
		},
		{
			name:      "Not Reach Limit",
			respCodes: []int{200, 200, 200, 200, 200, 200, 200, 200, 200, 200, http.StatusTooManyRequests},
		},
	}

	for _, test := range testCases {

		t.Run(test.name, func(t *testing.T) {
			// setup
			r := chi.NewRouter()
			handlerConfig := &Config{Router: r, ApiConfig: apiConfig}
			NewHandler(handlerConfig)

			for _, code := range test.respCodes {
				req := httptest.NewRequest("GET", "/", nil)
				rec := httptest.NewRecorder()
				r.ServeHTTP(rec, req)
				assert.Equal(t, code, rec.Result().StatusCode)

				var expectedResponseBody string
				if rec.Result().StatusCode == http.StatusOK {
					expectedResponseBody = "{\"status\":\"ok\"}"
				} else if rec.Result().StatusCode == http.StatusTooManyRequests {
					expectedResponseBody = "{\"type\":\"TOO_MANY_REQUESTS\",\"code\":\"api/too-many-requests\",\"message\":\"You are sending too many requests to the server.\",\"action\":\"C'mon buddy, please slow down. You are limited to 1 request/second.\"}"
				}

				assert.Equal(t, expectedResponseBody, strings.Trim(rec.Body.String(), "\n"))

			}

		})
	}
}
