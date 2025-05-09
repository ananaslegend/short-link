package clickhouse

import (
	"context"
	"fmt"

	"github.com/ananaslegend/short-link/internal/statistic/domain"
)

func (r StatisticRepository) AddRedirectEvent(
	ctx context.Context,
	redirectEvent domain.RedirectEventStatistic,
) error {
	const op = "internal.statistic.repository.clickhouse.StatisticRepository.AddRedirectEvent"

	err := r.conn.AsyncInsert(
		ctx,
		`insert into redirect_events (timestamp, link, alias)values (now(),?,?)`,
		false,
		redirectEvent.Link, redirectEvent.Alias,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
