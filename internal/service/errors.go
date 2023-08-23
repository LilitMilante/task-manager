package service

import (
	"errors"
)

var (
	ErrNotFound     = errors.New("not found")
	ErrUnauthorized = errors.New("unauthorized")
	ErrAccessDenied = errors.New("access denied")
)
