package public

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/siyoga/rollstory/pkg/db/postgres"
)

type Connection struct {
	*sqlx.DB
}

func New() (*Connection, error) {
	db, err := postgres.Connect()
	if err != nil {
		return nil, fmt.Errorf("while connecting to postgres: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("can't ping on fresh connection: %w", err)
	}

	return &Connection{db}, nil
}
