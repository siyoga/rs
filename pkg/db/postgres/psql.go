package postgres

import (
	"database/sql"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

func Connect() (*sqlx.DB, error) {
	dsnProvider := NewDSNProvider()

	// мб провайдить env снаружи, а не делать дефолт значения
	connConfig, options, err := buildConnectionConfigWithOptions(dsnProvider)
	if err != nil {
		return nil, err
	}

	return initConnect(connConfig, options)
}

func initConnect(connConfig *pgx.ConnConfig, options *Options) (*sqlx.DB, error) {
	db, err := sql.Open(options.driverName, stdlib.RegisterConnConfig(connConfig))
	if err != nil {
		return nil, err
	}

	return sqlx.NewDb(db, options.driverName), nil
}
