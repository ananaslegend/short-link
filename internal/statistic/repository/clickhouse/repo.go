package clickhouse

import (
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

type StatisticRepository struct {
	conn driver.Conn

	tracer trace.Tracer
}

func NewRepository(conn driver.Conn, traceProvider *sdktrace.TracerProvider) *StatisticRepository {
	return &StatisticRepository{conn: conn, tracer: traceProvider.Tracer("internal.statistic.repository.clickhouse.StatisticRepository")}
}
