package service

import (
	"context"

	"github.com/Jaytpa01/url-shortener-api/internal/entity"
)

// UrlService defines the methods the handler layer
// expects any url services it interacts with to implement.
type UrlService interface {
	ShortenUrl(ctx context.Context, url string) (*entity.Url, error)
	LengthenUrl(ctx context.Context, url string) (*entity.Url, error)
	GetUrlByToken(ctx context.Context, token string) (*entity.Url, error)
	IncrementUrlVisits(ctx context.Context, url *entity.Url) error
}
