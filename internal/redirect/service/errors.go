package service

import "errors"

var (
	ErrAliasNotFound = errors.New("alias not found")
	ErrEmptyAlias    = errors.New("alias is empty")
)
