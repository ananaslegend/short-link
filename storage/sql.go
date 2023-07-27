package storage

import (
	"database/sql"
	"fmt"
	"github.com/mattn/go-sqlite3"
)

type Sql struct {
	db *sql.DB
}

func (s Sql) PrepareStorage() error {
	const op = "storage.sql.PrepareStorage"

	stmt, err := s.db.Prepare(`
	create table if not exists link(
	    id integer primary key,
	    alias text not null unique,
	    link text not null);
	create index if not exists idx_alias on link(alias);
`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if _, err = stmt.Exec(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s Sql) AddLink(link, alias string) error {
	const op = "storage.sql.AddLink"

	stmt, err := s.db.Prepare(`
	insert into link (alias, link)
	values (?,?)
`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if _, err = stmt.Exec(link, alias); err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique { // TODO: sqlite3.Error прив'язались
			return fmt.Errorf("%s: \"%s\" %w", op, alias, ErrAliasExists)
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
