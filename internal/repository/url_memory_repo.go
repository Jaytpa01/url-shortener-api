package repository

import (
	"context"
	"sync"

	"github.com/Jaytpa01/url-shortener-api/internal/entity"
)

type memoryRepo struct {
	urls map[string]*entity.Url
	mu   sync.RWMutex
}

func NewInMemoryRepo() UrlRepository {
	return &memoryRepo{
		urls: make(map[string]*entity.Url),
		mu:   sync.RWMutex{},
	}
}

// Create is an in memory implementation of UrlRepository.Create
func (r *memoryRepo) Create(ctx context.Context, url *entity.Url) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.urls[url.Token]; exists {
		return ErrTokenAlreadyExists
	}

	r.urls[url.Token] = url
	return nil
}

// FindByToken is an in memory implementation of UrlRepository.FindByToken
func (r *memoryRepo) FindByToken(ctx context.Context, token string) (*entity.Url, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	url, ok := r.urls[token]
	if !ok {
		return nil, ErrUrlNotFound
	}

	return url, nil
}

// Update is an in memory implementation of UrlRepository.Update
func (r *memoryRepo) Update(ctx context.Context, url *entity.Url) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.urls[url.Token]
	if !ok {
		return ErrUrlNotFound
	}

	r.urls[url.Token] = url

	return nil
}

func (r *memoryRepo) GetAllUrls(ctx context.Context) ([]entity.Url, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	idx := 0
	urls := make([]entity.Url, len(r.urls))
	for _, url := range r.urls {
		urls[idx] = *url
		idx++
	}

	return urls, nil
}
