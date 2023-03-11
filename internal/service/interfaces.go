package service

import (
	"context"

	"github.com/Jaytpa01/url-shortener-api/internal/entity"
)

// UrlService defines the methods the handler layer
// expects any url services it interacts with to implement.
type UrlService interface {
	ShortenUrl(ctx context.Context, url string) (*entity.URL, error)
	LengthenUrl(ctx context.Context, url string) (*entity.URL, error)
	GetUrlByToken(ctx context.Context, token string) (*entity.URL, error)
}
