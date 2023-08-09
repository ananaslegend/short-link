package storage

import (
	_ "github.com/mattn/go-sqlite3"

	"database/sql"
	"fmt"
)

func NewSqliteStorage(dbPath string) (*Sql, error) {
	const op = "storage.sqlite.NewSqliteStorage"

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Sql{db: db}, nil
}
