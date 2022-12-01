package store

import "errors"

var (
	ErrNotFound   = errors.New("artisle not found")
	ErrNotContent = errors.New("not content")
)
