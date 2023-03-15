package mocks

import (
	"context"

	"github.com/Jaytpa01/url-shortener-api/internal/entity"
	"github.com/stretchr/testify/mock"
)

// mockUrlRepository is a mock implementation of our service.UrlRepository
type mockUrlRepository struct {
	mock.Mock
}

// NewMockUrlRepository returns a mock implementation of our UrlRepository for testing purposes.
// It is built using testify.Mock
func NewMockUrlRepository() *mockUrlRepository {
	return new(mockUrlRepository)
}

// FindByToken is a mock implementation of repository.FindByToken
func (m *mockUrlRepository) FindByToken(ctx context.Context, token string) (*entity.Url, error) {
	ret := m.Called(ctx, token)

	var r0 *entity.Url
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*entity.Url)
	}

	return r0, ret.Error(1)
}

// Create is a mock implementation of repository.Create
func (m *mockUrlRepository) Create(ctx context.Context, url *entity.Url) error {
	ret := m.Called(ctx, url)
	return ret.Error(0)
}

// Update is a mock implementation of repository.Update
func (m *mockUrlRepository) Update(ctx context.Context, url *entity.Url) error {
	ret := m.Called(ctx, url)
	return ret.Error(0)
}

// GetAllUrls is a mock implementation of repository.GetAllUrls
func (m *mockUrlRepository) GetAllUrls(ctx context.Context) ([]entity.Url, error) {
	ret := m.Called(ctx)

	var r0 []entity.Url
	if ret.Get(0) != nil {
		r0 = ret.Get(0).([]entity.Url)
	}

	return r0, ret.Error(1)
}
