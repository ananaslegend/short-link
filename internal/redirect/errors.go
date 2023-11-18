package redirect

import "errors"

var (
	ErrEmptyAlias     = errors.New("alias is empty")
	ErrAliasNotFound  = errors.New("alias not found")
	ErrCantSetToCache = errors.New("can't set to cache")
)
