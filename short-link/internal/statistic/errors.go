package statistic

import "errors"

var (
	ErrNoStatToInsert = errors.New("no stat to insert")
	ErrStatInserting  = errors.New("cant insert stat rows to db")
)
