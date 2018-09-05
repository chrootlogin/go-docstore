package database

import "errors"

var (
	ErrExists     = errors.New("object already exists")
	ErrNotFound   = errors.New("object not found")
	ErrNoFilename = errors.New("no filename")
)
