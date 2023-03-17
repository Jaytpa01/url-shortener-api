package repository

import (
	"context"
	"fmt"

	"github.com/Jaytpa01/url-shortener-api/internal/entity"
	"github.com/jmoiron/sqlx"
)

type sqliteRepository struct {
	db *sqlx.DB
}

func NewSQLiteRepository() (UrlRepository, error) {
	db, err := sqlx.Connect("sqlite3", "./db/url.db")
	if err != nil {
		return nil, fmt.Errorf("couldn't connect to sqlite database: %w", err)
	}

	repo := &sqliteRepository{
		db: db,
	}

	return repo, nil
}

func (s *sqliteRepository) FindByToken(ctx context.Context, token string) (*entity.Url, error) {
	panic("unimplemented") // TODO: implement
}

func (s *sqliteRepository) Create(ctx context.Context, url *entity.Url) error {
	panic("unimplemented") // TODO: implement
}

func (s *sqliteRepository) Update(ctx context.Context, url *entity.Url) error {
	panic("unimplemented") // TODO: implement
}

func (s *sqliteRepository) GetAllUrls(ctx context.Context) ([]entity.Url, error) {
	panic("unimplemented") // TODO: implement
}
