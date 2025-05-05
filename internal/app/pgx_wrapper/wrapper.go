package pgx_wrapper

import "github.com/jackc/pgx/v5/pgxpool"

type Wrapper struct {
	Pool *pgxpool.Pool
}
