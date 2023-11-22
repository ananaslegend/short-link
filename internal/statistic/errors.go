package statistic

import "errors"

var (
	ErrNoSatatToInsert = errors.New("no stat to insert")
	ErrCantInsertStat  = errors.New("cant insert stat rows to db")
)
