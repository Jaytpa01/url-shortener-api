package repository

import (
	"context"

	"github.com/Jaytpa01/url-shortener-api/internal/entity"
)

// UrlRepository defines the methods the service layer expects
// a url repository to implement.
type UrlRepository interface {
	FindByToken(ctx context.Context, token string) (*entity.Url, error)
	Create(ctx context.Context, url *entity.Url) error
	Update(ctx context.Context, url *entity.Url) error
	GetAllUrls(ctx context.Context) ([]entity.Url, error)
}
