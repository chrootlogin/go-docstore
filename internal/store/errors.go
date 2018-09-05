package store

import "errors"

var (
	ErrUserNotExist = errors.New("user does not exist")
	PointerIsNil    = errors.New("pointer is nil")
)
