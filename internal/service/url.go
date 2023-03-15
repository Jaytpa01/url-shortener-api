package service

import (
	"context"
	"errors"
	"fmt"

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

type Config struct {
	Logger  logger.Logger
	UrlRepo repository.UrlRepository
	Random  utils.Random
}

// urlService is used for the actual service implementation of this api
type urlService struct {
	logger  logger.Logger
	urlRepo repository.UrlRepository
	random  utils.Random
}

func NewUrlService(c *Config) UrlService {
	if c.Random == nil {
		c.Random = utils.NewRandomiser()
	}

	return &urlService{
		urlRepo: c.UrlRepo,
		logger:  c.Logger,
		random:  c.Random,
	}
}

// ShortenUrl validates the url, then attempts to create it in the repository.
func (u *urlService) ShortenUrl(ctx context.Context, url string) (*entity.Url, error) {
	if !validation.IsValidUrl(url) {
		return nil, api.NewBadRequest("url/invalid", fmt.Sprintf("The provided URL (%s) is invalid.", url))
	}

	newUrl := &entity.Url{
		Token:     u.random.GenerateRandomString(TOKEN_LENGTH),
		TargetUrl: url,
	}

	var err error
	// because there's a very slim chance a generated token may clash, we give it up to 3 attempts
	attempts := 3
	for i := 0; i < attempts; i++ {
		err = u.urlRepo.Create(ctx, newUrl)
		if err != nil {
			if errors.Is(err, repository.ErrTokenAlreadyExists) {
				newUrl.Token = u.random.GenerateRandomString(TOKEN_LENGTH)
				continue
			} else {
				return nil, api.NewInternal("couldnt-shorten", api.WithDebug(err.Error()))
			}
		}

		// no error, so we break
		break
	}

	if err != nil {
		// Because we try 3 times to prevent a token clash, we simply return an internal error here after several
		// attempts if we keep clashinig
		apiErr := api.NewInternal("couldnt-shorten", api.WithDebug(err.Error()))
		u.logger.Info("Failed to shorten URL.", apiErr)
		return nil, apiErr
	}

	return newUrl, nil
}

// LengthenUrl validates the url, then attempts to create it in the repository.
func (u *urlService) LengthenUrl(ctx context.Context, url string) (*entity.Url, error) {
	if !validation.IsValidUrl(url) {
		return nil, api.NewBadRequest("url/invalid", "The provided URL (%s) is invalid.")
	}

	newUrl := &entity.Url{
		Token:     u.random.GenerateRandomString(utils.Max(MINIMUM_LONG_TOKEN_LENGTH, len(url)*LENGTHEN_TOKEN_SCALE_FACTOR)),
		TargetUrl: url,
	}

	var err error
	// because there's a very slim chance a generated token may clash, we give it up to 3 attempts
	for i := 0; i < 3; i++ {
		err = u.urlRepo.Create(ctx, newUrl)
		if err != nil {
			if errors.Is(err, repository.ErrTokenAlreadyExists) {
				newUrl.Token = u.random.GenerateRandomString(utils.Max(MINIMUM_LONG_TOKEN_LENGTH, len(url)*LENGTHEN_TOKEN_SCALE_FACTOR))
				continue
			} else {
				return nil, err
			}
		}

		// no error, so we break
		break
	}

	if err != nil {
		// Because createLongUrl recursively retries until theres no token clash, we simply return
		// an internal error here

		apiErr := api.NewInternal("couldnt-lengthen", api.WithDebug(err.Error()))
		u.logger.Info("Failed to lengthen URL.", apiErr)
		return nil, apiErr
	}

	return newUrl, nil
}

// FindUrlByToken attempts to find the URL associated to provided token.
func (u *urlService) FindUrlByToken(ctx context.Context, token string) (*entity.Url, error) {
	url, err := u.urlRepo.FindByToken(ctx, token)
	if err != nil {
		if errors.Is(err, repository.ErrUrlNotFound) {
			return nil, api.NewNotFound("url/not-found", fmt.Sprintf("Couldn't find URL with token (%s).", token))
		}

		apiErr := api.NewInternal("url/internal", api.WithDebug(err.Error()))
		u.logger.Info("Couldn't retrieve URL", apiErr)

		return nil, apiErr
	}

	return url, nil
}

func (u *urlService) IncrementUrlVisits(ctx context.Context, url *entity.Url) error {
	url.Visits++
	return u.urlRepo.Update(ctx, url)
}

func (u *urlService) GetAllUrls(ctx context.Context) ([]entity.Url, error) {
	return u.urlRepo.GetAllUrls(ctx)
}
