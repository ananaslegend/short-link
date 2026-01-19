package domain

import "time"

type RedirectEventStatistic struct {
	LinkID     int64
	Link       string
	Alias      string
	OccurredAt time.Time
}
