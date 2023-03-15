package repository

import "errors"

var (
	ErrUrlNotFound        = errors.New("url not found")
	ErrTokenAlreadyExists = errors.New("token already exists")
)
