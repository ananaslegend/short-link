package clickhouse

import (
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type StatisticRepository struct {
	conn driver.Conn
}

func NewRepository(conn driver.Conn, traceProvider *sdktrace.TracerProvider) *StatisticRepository {
	return &StatisticRepository{
		conn: conn,
	}
}
