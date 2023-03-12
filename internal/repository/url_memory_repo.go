package repository

import (
	"context"
	"errors"
	"sync"

	"github.com/Jaytpa01/url-shortener-api/internal/entity"
)

var (
	ErrTokenNotFound  = errors.New("token not found")
	ErrUrlDoesntExist = errors.New("url doesn't exist")
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

func (r *memoryRepo) FindByToken(ctx context.Context, token string) (*entity.Url, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	url, ok := r.urls[token]
	if !ok {
		return nil, ErrTokenNotFound
	}

	return url, nil
}

func (r *memoryRepo) GetAll(ctx context.Context) ([]*entity.Url, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	idx := 0
	urls := make([]*entity.Url, len(r.urls))
	for _, url := range r.urls {
		urls[idx] = url
		idx++
	}

	return urls, nil
}

func (r *memoryRepo) Create(ctx context.Context, url *entity.Url) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.urls[url.Token] = url
	return nil
}

func (r *memoryRepo) Update(ctx context.Context, url *entity.Url) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.urls[url.Token]
	if !ok {
		return ErrUrlDoesntExist
	}

	r.urls[url.Token] = url

	return nil
}
