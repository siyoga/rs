package public

import (
	"context"
	"database/sql"
)

// DB Реализуется через *sqlx.DB
type DB interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	Begin() (*sql.Tx, error)
}
