package storage

import "errors"

var (
	ErrAliasNotFound = errors.New("alias not found")
	ErrAliasExists   = errors.New("alias exists")
)
