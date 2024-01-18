package repository

import "errors"

var (
	ErrAliasNotFound  = errors.New("alias not found")
	ErrCantSetToCache = errors.New("can't set to cache")
)
