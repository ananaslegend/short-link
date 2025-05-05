package postgres

import "errors"

var (
	ErrAliasNtFound       = errors.New("alias not found")
	ErrAliasAlreadyExists = errors.New("alias already exists")
)
