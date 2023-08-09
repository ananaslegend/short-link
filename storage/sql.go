package storage

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/ananaslegend/short-link/errs"
	"github.com/mattn/go-sqlite3"
)

type Sql struct {
	db *sql.DB
}

func (s *Sql) Close() {
	err := s.db.Close()
	if err != nil {
		return
	}
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

func (s Sql) InsertLink(link, alias string) error {
	const op = "storage.sql.InsertLink"

	stmt, err := s.db.Prepare(`
	insert into link (alias, link)
	values (?,?)
`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	defer stmt.Close()

	if _, err := stmt.Exec(alias, link); err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique { // TODO: sqlite3.Error прив'язались
			return fmt.Errorf("%s: \"%s\" %w", op, alias, errs.ErrAliasExists)
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s Sql) SelectLink(ctx context.Context, alias string) (string, error) {
	const op = "storage.sql.SelectLink"

	stmt, err := s.db.Prepare(`
	select link 
	from link
	where alias == ?
`)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	var link string
	if err = stmt.QueryRowContext(ctx, alias).Scan(&link); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return link, nil
}
