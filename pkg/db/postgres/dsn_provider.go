package postgres

import (
	"errors"
	"fmt"
	"github.com/siyoga/rollstory/pkg/env"
	"os"
)

type DSNProvider struct {
}

func NewDSNProvider() DSNProvider {
	return DSNProvider{}
}

func (p *DSNProvider) IsAvailable(databaseENV databaseENV) bool {
	_, err := p.readENV(databaseENV)

	return err == nil
}

func (p *DSNProvider) Provide(databaseENV databaseENV) (*DSN, error) {
	parsedDSN, err := p.readENV(databaseENV)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ENV: %w", err)
	}

	dsn, err := p.buildDSN(parsedDSN, int(UserRoleFullAccess))
	if err != nil {
		return nil, fmt.Errorf("failed to build DSN: %w", err)
	}

	return dsn, nil
}

func (p *DSNProvider) readENV(databaseENV databaseENV) (databaseDSN, error) {
	errs := make([]error, 0, 7) // 7 env vars

	host, err := env.LookupEnvString(databaseENV.Host)
	if err != nil {
		errs = append(errs, fmt.Errorf("environment variable %s not set", databaseENV.Host))
	}

	port, err := env.LookupEnvInt(databaseENV.Port)
	if err != nil {
		errs = append(errs, fmt.Errorf("environment variable %s not set", databaseENV.Port))
	}

	dbName, err := env.LookupEnvString(databaseENV.DbName)
	if err != nil {
		errs = append(errs, fmt.Errorf("environment variable %s not set", databaseENV.DbName))
	}

	user, err := env.LookupEnvString(databaseENV.User)
	if err != nil {
		errs = append(errs, fmt.Errorf("environment variable %s not set", databaseENV.User))
	}

	password, err := env.LookupEnvString(databaseENV.Password)
	if err != nil {
		errs = append(errs, fmt.Errorf("environment variable %s not set", databaseENV.Password))
	}

	sslMode, err := env.LookupEnvString(databaseENV.SSLMode)
	if err != nil {
		errs = append(errs, fmt.Errorf("environment variable %s not set", databaseENV.SSLMode))
	}

	connectTimeout, err := env.LookupEnvInt(databaseENV.ConnectTimeout)
	if err != nil {
		errs = append(errs, fmt.Errorf("environment variable %s not set", databaseENV.ConnectTimeout))
	}

	if len(errs) != 0 {
		return databaseDSN{}, errors.Join(errs...)
	}

	return databaseDSN{
		Name:           "",
		Host:           host,
		Port:           port,
		DbName:         dbName,
		SSLMode:        sslMode,
		ConnectTimeout: connectTimeout,
		User: databaseUsers{
			FullAccess: &databaseUser{
				Login:    user,
				Password: password,
			},
		},
	}, nil
}

func (p *DSNProvider) buildDSN(config databaseDSN, userRole int) (*DSN, error) {
	// пока один дефолтный юзер на full access права
	var u *databaseUser
	switch userRole {
	case int(UserRoleFullAccess):
		fallthrough
	default:
		u = config.User.FullAccess
	}

	if u == nil {
		return nil, fmt.Errorf("invalid user role %s", UserRole(userRole))
	}

	userLogin := u.Login
	if envVal := os.ExpandEnv(userLogin); envVal != "" {
		userLogin = envVal
	}

	userPassword := u.Password
	if envVal := os.ExpandEnv(userPassword); envVal != "" {
		userPassword = envVal
	}

	dsn := DSN{
		Name:           config.Name,
		DBName:         config.Name,
		Host:           config.Host,
		Port:           config.Port,
		ConnectTimeout: config.ConnectTimeout,
		UserLogin:      userLogin,
		UserPassword:   userPassword,
		UserRole:       userRole,
		SSLMode:        config.SSLMode,
	}

	return &dsn, nil
}
