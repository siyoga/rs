package postgres

import (
	"fmt"
	"github.com/jackc/pgx/v4"
	"net/url"
)

func buildOptions(dsnProvider DSNProvider) (*Options, error) {
	options := &Options{}

	env := databaseENV{
		Host:           "PG_HOST",
		Port:           "PG_PORT",
		DbName:         "PG_DBNAME",
		User:           "PG_USER",
		Password:       "PG_PASSWORD",
		SSLMode:        "PG_SSLMODE",
		ConnectTimeout: "PG_CONNECT_TIMEOUT",
	}

	if dsnProvider.IsAvailable(env) {
		dsn, err := dsnProvider.Provide(env)
		if err != nil {
			return nil, fmt.Errorf("when providing dsn: %w", err)
		}

		if err := parseDSNToOptions(dsn, options); err != nil {
			return nil, fmt.Errorf("when parsing dsn: %w", err)
		}
	}

	return options, nil
}

func buildConnectionConfigWithOptions(dsnProvider DSNProvider) (*pgx.ConnConfig, *Options, error) {
	options, err := buildOptions(dsnProvider)
	if err != nil {
		return nil, nil, fmt.Errorf("when building connection config: %w", err)
	}

	connConfig, err := pgx.ParseConfig(buildDSNUrlFromOptions(options))
	if err != nil {
		return nil, nil, fmt.Errorf("when building connection config: %w", err)
	}

	useOptions(options, connConfig)

	return connConfig, options, nil
}

func useOptions(o *Options, config *pgx.ConnConfig) {
	if o.host != "" {
		config.Host = o.host
	}

	if o.port != 0 {
		config.Port = uint16(o.port)
	}

	if o.user != "" {
		config.User = o.user
	}

	if o.password != "" {
		config.Password = o.password
	}

	if o.dbName != "" {
		config.Database = o.dbName
	}

	if o.Logger != nil {
		config.Logger = o.Logger
	}

	if o.LogLevel != 0 {
		config.LogLevel = o.LogLevel
	}
}

func parseDSNToOptions(dsn *DSN, options *Options) error {
	if dsn == nil {
		return fmt.Errorf("dsn is nil")
	}

	options.user = dsn.UserLogin
	options.password = dsn.UserPassword
	options.dsnName = dsn.Name
	options.host = dsn.Host
	options.port = dsn.Port
	options.dbName = dsn.DBName
	options.userRole = UserRole(dsn.UserRole)
	options.sslMode = SSLModeDisabled
	options.connectTimeout = dsn.ConnectTimeout

	return nil
}

func buildDSNUrlFromOptions(options *Options) string {
	if options.host == "" || options.port == 0 {
		return ""
	}

	dbURL := url.URL{
		Scheme: dsnSchemePostgreSQL,
		User:   url.UserPassword(options.user, options.password),
		Host:   fmt.Sprintf("%s:%d", options.host, options.port),
		Path:   options.dbName,
	}

	q := dbURL.Query()
	q.Add("sslmode", string(options.sslMode))

	if options.connectTimeout > 0 {
		q.Add("connect_timeout", fmt.Sprintf("%d", options.connectTimeout))
	}

	dbURL.RawQuery = q.Encode()

	return dbURL.String()
}
