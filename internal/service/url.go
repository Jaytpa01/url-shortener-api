package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Jaytpa01/url-shortener-api/api"
	"github.com/Jaytpa01/url-shortener-api/internal/entity"
	"github.com/Jaytpa01/url-shortener-api/internal/repository"
	"github.com/Jaytpa01/url-shortener-api/pkg/logger"
	"github.com/Jaytpa01/url-shortener-api/pkg/utils"
	"github.com/Jaytpa01/url-shortener-api/pkg/validation"
)

const (
	// TODO: Inject this value via a config
	TOKEN_LENGTH = 6

	LENGTHEN_TOKEN_SCALE_FACTOR = 2  // double the length of URLs
	MINIMUM_LONG_TOKEN_LENGTH   = 42 // the mimimum length that a lengthened token should be
)

type urlService struct {
	logger  logger.Logger
	urlRepo repository.UrlRepository
}

type Config struct {
	Logger  logger.Logger
	UrlRepo repository.UrlRepository
}

func NewUrlService(c *Config) UrlService {
	return &urlService{
		urlRepo: c.UrlRepo,
		logger:  c.Logger,
	}
}

// ShortenUrl validates the url, then attempts to create it in the repository.
func (u *urlService) ShortenUrl(ctx context.Context, url string) (*entity.Url, error) {
	if !validation.IsValidUrl(url) {
		return nil, api.NewBadRequest("invalid-url", "The provided URL (%s) is invalid.")
	}

	newUrl := &entity.Url{
		Token:     utils.GenerateRandomString(TOKEN_LENGTH),
		TargetUrl: url,
		CreatedAt: time.Now(),
	}

	err := createShortUrl(ctx, newUrl, u.urlRepo)
	if err != nil {
		// Because createShortUrl recursively retries until theres no token clash, we simply return
		// an internal error here

		apiErr := api.NewInternal("couldnt-shorten", api.WithDebug(err.Error()))
		u.logger.Info("Failed to shorten URL.", apiErr)
		return nil, apiErr
	}

	return newUrl, nil
}

// LengthenUrl validates the url, then attempts to create it in the repository.
func (u *urlService) LengthenUrl(ctx context.Context, url string) (*entity.Url, error) {
	if !validation.IsValidUrl(url) {
		return nil, api.NewBadRequest("invalid-url", "The provided URL (%s) is invalid.")
	}

	newUrl := &entity.Url{
		Token:     utils.GenerateRandomString(utils.Max(MINIMUM_LONG_TOKEN_LENGTH, len(url)*LENGTHEN_TOKEN_SCALE_FACTOR)),
		TargetUrl: url,
		CreatedAt: time.Now(),
	}

	err := createLongUrl(ctx, newUrl, u.urlRepo)
	if err != nil {
		// Because createLongUrl recursively retries until theres no token clash, we simply return
		// an internal error here

		apiErr := api.NewInternal("couldnt-lengthen", api.WithDebug(err.Error()))
		u.logger.Info("Failed to lengthen URL.", apiErr)
		return nil, apiErr
	}

	return newUrl, nil
}

// GetUrlByToken attempts to find the URL associated to provided token.
func (u *urlService) GetUrlByToken(ctx context.Context, token string) (*entity.Url, error) {
	url, err := u.urlRepo.FindByToken(ctx, token)
	if err != nil {
		if errors.Is(err, api.ErrUrlNotFound) {
			return nil, api.NewNotFound("url-not-found", fmt.Sprintf("Couldn't find URL with token (%s).", token))
		}

		apiErr := api.NewInternal("internal-error-url-not-found", api.WithDebug(err.Error()))
		u.logger.Info("Couldn't retrieve URL", apiErr)

		return nil, apiErr
	}

	return url, nil
}

func (u *urlService) IncrementUrlVisits(ctx context.Context, url *entity.Url) error {
	url.Visits++
	return u.urlRepo.Update(ctx, url)
}

// createShortUrl attempts to create a new url resource. If the token already exists,
// it will recursively try with a new token. If its some other error, it will simply return it.
func createShortUrl(ctx context.Context, url *entity.Url, repo repository.UrlRepository) error {
	err := repo.Create(ctx, url)
	if errors.Is(err, api.ErrTokenAlreadyExists) {
		url.Token = utils.GenerateRandomString(TOKEN_LENGTH)
		return createShortUrl(ctx, url, repo)
	}
	return err
}

// createLongUrl attempts to create a new url resource. If the token already exists,
// it will recursively try with a new token. If its some other error, it will simply return it.
func createLongUrl(ctx context.Context, url *entity.Url, repo repository.UrlRepository) error {
	err := repo.Create(ctx, url)
	if errors.Is(err, api.ErrTokenAlreadyExists) {
		url.Token = utils.GenerateRandomString(utils.Max(MINIMUM_LONG_TOKEN_LENGTH, len(url.TargetUrl)*LENGTHEN_TOKEN_SCALE_FACTOR))
		return createShortUrl(ctx, url, repo)
	}
	return err
}
