package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Jaytpa01/url-shortener-api/internal/entity"
	"github.com/jmoiron/sqlx"
)

type sqliteRepository struct {
	db *sqlx.DB
}

func NewSQLiteRepository(dsn string) (UrlRepository, error) {
	db, err := sqlx.Connect("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("couldn't connect to sqlite database: %w", err)
	}

	repo := &sqliteRepository{
		db: db,
	}

	return repo, nil
}

func (s *sqliteRepository) FindByToken(ctx context.Context, token string) (*entity.Url, error) {
	url := &entity.Url{}

	err := s.db.GetContext(ctx, url, "SELECT * FROM url WHERE token = ?", token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUrlNotFound
		}

		return nil, err
	}

	return url, nil
}

func (s *sqliteRepository) Create(ctx context.Context, url *entity.Url) error {
	// Create a helper function for preparing failure results.
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	// Defer a rollback in case anything fails.
	defer tx.Rollback()

	result, err := tx.ExecContext(ctx, `INSERT INTO url (token, target_url, created_at) VALUES (?, ?, ?)`, url.Token, url.TargetUrl, url.CreatedAt)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected != 1 {
		return fmt.Errorf("%d rows affected, expected 1", rowsAffected)
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *sqliteRepository) Update(ctx context.Context, url *entity.Url) error {
	panic("unimplemented") // TODO: implement
}

func (s *sqliteRepository) GetAllUrls(ctx context.Context) ([]entity.Url, error) {
	panic("unimplemented") // TODO: implement
}
