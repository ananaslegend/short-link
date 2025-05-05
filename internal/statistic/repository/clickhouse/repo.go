package clickhouse

import "github.com/ClickHouse/clickhouse-go/v2/lib/driver"

type Statistic struct {
	conn driver.Conn
}

func NewRepository(conn driver.Conn) *Statistic {
	return &Statistic{conn: conn}
}
