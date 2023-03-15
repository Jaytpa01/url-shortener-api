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
			&entity.Url{
				Token:     "123456",
				TargetUrl: "https://example.com",
			},
			repository.ErrTokenAlreadyExists,
			api.NewInternal("couldnt-shorten", api.WithDebug(repository.ErrTokenAlreadyExists.Error())),
		},
		{
			"Unknown Repo Error",
			"https://example.com",
			&entity.Url{
				Token:     "123456",
				TargetUrl: "https://example.com",
			},
			errors.New("some repo error"),
			api.NewInternal("couldnt-shorten", api.WithDebug("some repo error")),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := mocks.NewMockUrlRepository()
			repo.On("Create", mock.AnythingOfType("*context.emptyCtx"), test.url).Return(test.repoError)

			randomiser := mocks.NewMockRandomiser()
			randomiser.On("GenerateRandomString", 6).Return("123456")
			urlService := NewUrlService(&Config{
				UrlRepo: repo,
				Logger:  logger.NewApiLogger(),
				Random:  randomiser,
			})

			url, err := urlService.ShortenUrl(context.Background(), test.input)

			if test.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, test.expectedError, err)
			}

			if test.url != nil && test.repoError == nil {
				assert.Equal(t, test.url, url)
			}

		})
	}
}

// TODO: Write tests for LengthenUrl

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
			Logger:  logger.NewApiLogger(),
		})

		result, err := urlService.FindUrlByToken(context.Background(), test.token)
		if test.expectedError != nil {
			assert.Error(t, err)
			assert.Equal(t, test.expectedError, err)
		}

		assert.Equal(t, test.returnUrl, result)

	}

}

// TODO: Write tests for IncrementUrlVisits
