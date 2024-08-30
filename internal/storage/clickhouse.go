package storage

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"time"
)

func NewClickHouseDB(host, port, db, pass, user string) (*sql.DB, error) {
	chDB := clickhouse.OpenDB(&clickhouse.Options{
		Addr: []string{fmt.Sprintf("%s:%s", host, port)},
		Auth: clickhouse.Auth{
			Database: db,
			Username: user,
			Password: pass,
		},
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
	})

	if err := chDB.Ping(); err != nil {
		return nil, err
	}

	return chDB, nil
}

func NewClickHouse(host, port, db, pass, user string) (driver.Conn, error) {
	chDB, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{fmt.Sprintf("%s:%s", host, port)},
		Auth: clickhouse.Auth{
			Database: db,
			Username: user,
			Password: pass,
		},
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
	})

	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = chDB.Ping(ctx); err != nil {
		return nil, err
	}

	return chDB, nil
}
