package service

import "errors"

var (
	ErrAliasAlreadyExists     = errors.New("alias exists")
	ErrAutoAliasAlreadyExists = errors.New("auto alias already exists")
)
