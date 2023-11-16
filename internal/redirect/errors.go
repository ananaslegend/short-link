package redirect

import "errors"

var (
	ErrEmptyAlias    = errors.New("alias is empty")
	ErrAliasNotFound = errors.New("alias not found")
)
