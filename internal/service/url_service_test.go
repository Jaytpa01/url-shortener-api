package service

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/Jaytpa01/url-shortener-api/api"
	"github.com/Jaytpa01/url-shortener-api/internal/entity"
	"github.com/Jaytpa01/url-shortener-api/internal/mocks"
	"github.com/Jaytpa01/url-shortener-api/internal/repository"
	"github.com/Jaytpa01/url-shortener-api/pkg/logger"
	"github.com/Jaytpa01/url-shortener-api/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_ShortenUrl(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		url           *entity.Url
		repoError     error
		expectedError error
	}{
		{
			"Create Valid Url",
			"https://example.com",
			&entity.Url{
				Token:     "123456",
				TargetUrl: "https://example.com",
			},
			nil,
			nil,
		},
		{
			"Create Invalid Url",
			"example",
			nil,
			nil,
			api.NewBadRequest("url/invalid", fmt.Sprintf("The provided URL (%s) is invalid.", "example")),
		},
		{
			"No Url",
			"",
			nil,
			nil,
			api.NewBadRequest("url/invalid", fmt.Sprintf("The provided URL (%s) is invalid.", "")),
		},
		{
			"Token Clash",
			"https://example.com",
			nil,
			repository.ErrTokenAlreadyExists,
			api.NewInternal("url/couldnt-shorten", api.WithDebug(repository.ErrTokenAlreadyExists.Error())),
		},
		{
			"Unknown Repo Error",
			"https://example.com",
			&entity.Url{
				Token:     "123456",
				TargetUrl: "https://example.com",
			},
			errors.New("some repo error"),
			api.NewInternal("url/couldnt-shorten", api.WithDebug("some repo error")),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			repo := mocks.NewMockUrlRepository()
			repo.On("Create", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("*entity.Url")).Return(test.repoError)

			randomiser := mocks.NewMockRandomiser()
			randomiser.On("GenerateRandomString", 6).Return("123456")
			urlService := NewUrlService(&Config{
				UrlRepo: repo,
				Logger:  logger.NewApiLogger("development"),
				Random:  randomiser,
			})

			url, err := urlService.ShortenUrl(context.Background(), test.input)

			if test.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, test.expectedError, err)
			}

			if test.url != nil && test.repoError == nil {
				// Because the service implementation uses time.Now() for the CreatedAt timestamp, the test.url and url returned
				// will almost always be slightly different because of this. We could use an interface and use dependeny injection,
				// but thats unneccessary and annoying, so we just coerce the values to be equal here
				createdAt := time.Now()
				test.url.CreatedAt = createdAt
				url.CreatedAt = createdAt

				assert.Equal(t, test.url, url)
			}

		})
	}
}

func Test_LengthenUrl(t *testing.T) {
	testCases := []struct {
		name           string
		inputUrl       string
		serviceUrlResp *entity.Url
		repoError      error
		expectedError  error
	}{
		{

			"Successfully Lengthen URL",
			"https://example.com",
			&entity.Url{
				Token:     "ThisIsMeantToRepresentAReallyReallyReallyReallyLongToken",
				TargetUrl: "https://example.com",
			},
			nil,
			nil,
		},
		{
			"Create Invalid Url",
			"example",
			nil,
			repository.ErrTokenAlreadyExists,
			api.NewBadRequest("url/invalid", fmt.Sprintf("The provided URL (%s) is invalid.", "example")),
		},
		{
			"Some Unknown Repo Error",
			"https://example.com",
			nil,
			errors.New("some unknown error"),
			api.NewInternal("url/couldnt-lengthen", api.WithDebug(errors.New("some unknown error").Error())),
		},
		{
			"Token Clash",
			"https://example.com",
			nil,
			repository.ErrTokenAlreadyExists,
			api.NewInternal("url/couldnt-lengthen", api.WithDebug(repository.ErrTokenAlreadyExists.Error())),
		},
	}

	for _, test := range testCases {
		repo := mocks.NewMockUrlRepository()
		repo.On("Create", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("*entity.Url")).Return(test.repoError)

		randomiser := mocks.NewMockRandomiser()
		randomiser.On("GenerateRandomString", utils.Max(MINIMUM_LONG_TOKEN_LENGTH, len(test.inputUrl)*LENGTHEN_TOKEN_SCALE_FACTOR)).Return("ThisIsMeantToRepresentAReallyReallyReallyReallyLongToken")
		urlService := NewUrlService(&Config{
			UrlRepo: repo,
			Logger:  logger.NewApiLogger("development"),
			Random:  randomiser,
		})

		url, err := urlService.LengthenUrl(context.Background(), test.inputUrl)

		if test.expectedError != nil {
			assert.Error(t, err)
			assert.Equal(t, test.expectedError, err)
		}

		if test.serviceUrlResp != nil && test.repoError == nil {
			// Because the service implementation uses time.Now() for the CreatedAt timestamp, the test.serviceUrlResp and url returned
			// will almost always be slightly different because of this. We could use an interface and use dependeny injection,
			// but thats unneccessary and annoying, so we just coerce the values to be equal here
			createdAt := time.Now()
			test.serviceUrlResp.CreatedAt = createdAt
			url.CreatedAt = createdAt

			assert.Equal(t, test.serviceUrlResp, url)
		}
	}
}

func Test_FindUrlByToken(t *testing.T) {
	testCases := []struct {
		name          string
		token         string
		returnUrl     *entity.Url
		repoError     error
		expectedError error
	}{
		{
			"Url Exists With Given Token",
			"au42Gq",
			&entity.Url{
				Token:     "au42Gq",
				TargetUrl: "https://example.com",
				Visits:    42,
				CreatedAt: time.Now(),
			},
			nil,
			nil,
		},
		{
			"Url Does NOT Exist With Given Token",
			"qwerty",
			nil,
			repository.ErrUrlNotFound,
			api.NewNotFound("url/not-found", fmt.Sprintf("Couldn't find URL with token (%s).", "qwerty")),
		},
		{
			"Unknown Repository Error",
			"987654",
			nil,
			errors.New("whoops, some error"),
			api.NewInternal("url/internal", api.WithDebug("whoops, some error")),
		},
	}

	for _, test := range testCases {
		repo := mocks.NewMockUrlRepository()
		repo.On("FindByToken", mock.AnythingOfType("*context.emptyCtx"), test.token).Return(test.returnUrl, test.repoError)

		urlService := NewUrlService(&Config{
			UrlRepo: repo,
			Logger:  logger.NewApiLogger("development"),
		})

		result, err := urlService.FindUrlByToken(context.Background(), test.token)
		if test.expectedError != nil {
			assert.Error(t, err)
			assert.Equal(t, test.expectedError, err)
		}

		assert.Equal(t, test.returnUrl, result)

	}

}

func Test_IncrementUrlVisits(t *testing.T) {
	testCases := []struct {
		name                            string
		inputUrl                        *entity.Url
		expectedVisitsAfterIncrementing int
		repoError                       error
		expectedServiceError            error
	}{
		{
			"Successfully Incremented Visits",
			&entity.Url{
				Token:  "987654",
				Visits: 0,
			},
			1,
			nil,
			nil,
		},
		{
			"Couldn't Increment Visits",
			&entity.Url{
				Token:  "987654",
				Visits: 0,
			},
			0,
			errors.New("couldn't increment"),
			errors.New("couldn't increment"),
		},
		{
			"Provide Nil URL Pointer",
			nil,
			0,
			ErrNilUrlPointer,
			ErrNilUrlPointer,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			repo := mocks.NewMockUrlRepository()
			repo.On("Update", mock.AnythingOfType("*context.emptyCtx"), test.inputUrl).Return(test.repoError)

			service := NewUrlService(&Config{
				UrlRepo: repo,
				Logger:  logger.NewApiLogger("development"),
			})

			err := service.IncrementUrlVisits(context.Background(), test.inputUrl)
			if test.expectedServiceError != nil {
				assert.Error(t, err)
				assert.Equal(t, test.expectedServiceError, err)
			}

			if test.inputUrl != nil {
				assert.Equal(t, test.expectedVisitsAfterIncrementing, test.inputUrl.Visits)
			}
		})
	}
}

func Test_GetAllUrls(t *testing.T) {
	testCases := []struct {
		name            string
		repoUrlResponse []entity.Url
		repoErr         error
		expectedErr     error
	}{
		{
			"Success",
			[]entity.Url{
				{
					Token:     "123456",
					TargetUrl: "https://example.com",
				},
				{
					Token:     "987654",
					TargetUrl: "https://google.com",
				},
			},
			nil,
			nil,
		},
		{
			"Failure",
			nil,
			errors.New("some error"),
			errors.New("some error"),
		},
		{
			"Success But No URLS in Repo",
			[]entity.Url{},
			nil,
			nil,
		},
	}

	for _, test := range testCases {
		repo := mocks.NewMockUrlRepository()
		repo.On("GetAllUrls", mock.AnythingOfType("*context.emptyCtx")).Return(test.repoUrlResponse, test.repoErr)

		service := NewUrlService(&Config{
			UrlRepo: repo,
			Logger:  logger.NewApiLogger("development"),
		})

		urls, err := service.GetAllUrls(context.Background())

		if test.expectedErr != nil {
			assert.Error(t, err)
			assert.Equal(t, test.expectedErr, err)
		}

		assert.Equal(t, test.repoUrlResponse, urls)

	}

}
