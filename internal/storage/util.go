package storage

import (
	"database/sql"
	"fmt"
)

func Prepare(db *sql.DB) error {
	const op = "storage.sql.Prepare"

	stmt, err := db.Prepare(`
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

	stmt, err = db.Prepare(`
		create table if not exists statistic(
	    id integer primary key,
	    redirect_time_stamp integer,
	    link text,
		redirect integer                      
	);`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if _, err = stmt.Exec(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
