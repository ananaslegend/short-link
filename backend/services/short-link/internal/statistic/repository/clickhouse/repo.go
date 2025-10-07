package clickhouse

import (
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type StatisticRepository struct {
	conn driver.Conn
}

func NewRepository(conn driver.Conn) *StatisticRepository {
	return &StatisticRepository{
		conn: conn,
	}
}
