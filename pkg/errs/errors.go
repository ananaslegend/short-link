package errs

import "errors"

var (
	ErrAliasNotFound          = errors.New("alias not found")
	ErrAliasExists            = errors.New("alias exists")
	ErrAutoAliasAlreadyExists = errors.New("auto alias already exists")

	ErrEmptyAlias = errors.New("alias is empty")
)
