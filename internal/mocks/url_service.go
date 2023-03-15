package mocks

import (
	"context"

	"github.com/Jaytpa01/url-shortener-api/internal/entity"
	"github.com/stretchr/testify/mock"
)

// mockUrlService is a mock implementation of our service.UrlService
type mockUrlService struct {
	mock.Mock
}

// NewMockUrlService returns a mock implementation of our UrlService for testing purposes.
// It is built using testify.Mock
func NewMockUrlService() *mockUrlService {
	return new(mockUrlService)
}

// ShortenUrl is a mock implementation of UrlService.ShortenUrl
func (m *mockUrlService) ShortenUrl(ctx context.Context, url string) (*entity.Url, error) {
	ret := m.Called(ctx, url)

	var r0 *entity.Url
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*entity.Url)
	}

	return r0, ret.Error(1)
}

// LengthenUrl is a mock implementation of UrlService.LengthenUrl
func (m *mockUrlService) LengthenUrl(ctx context.Context, url string) (*entity.Url, error) {
	ret := m.Called(ctx, url)

	var r0 *entity.Url
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*entity.Url)
	}

	return r0, ret.Error(1)
}

// FindUrlByToken is a mock implementation of UrlService.FindUrlByToken
func (m *mockUrlService) FindUrlByToken(ctx context.Context, token string) (*entity.Url, error) {
	ret := m.Called(ctx, token)

	var r0 *entity.Url
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*entity.Url)
	}

	return r0, ret.Error(1)
}

// IncrementUrlVisits is a mock implementation of UrlService.IncrementUrlVisits
func (m *mockUrlService) IncrementUrlVisits(ctx context.Context, url *entity.Url) error {
	ret := m.Called(ctx, url)

	return ret.Error(0)
}

func (m *mockUrlService) GetAllUrls(ctx context.Context) ([]entity.Url, error) {
	ret := m.Called(ctx)

	var r0 []entity.Url
	if ret.Get(0) != nil {
		r0 = ret.Get(0).([]entity.Url)
	}

	return r0, ret.Error(1)
}
