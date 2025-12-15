package postgres

import (
	"fmt"
	"github.com/jackc/pgx/v4"
	"strings"
)

const (
	dsnSchemePostgreSQL = "postgresql"
)

type SSLModeType string

const (
	SSLModeDisabled SSLModeType = "disabled"
)

type UserRole int

const (
	UserRoleReadOnly = UserRole(iota)
	UserRoleReadWrite
	UserRoleFullAccess
)

type Connection struct {
	sqlx
}

type Options struct {
	host           string
	port           int
	user           string
	password       string
	dbName         string
	dsnName        string
	driverName     string
	connectTimeout int
	sslMode        SSLModeType
	userRole       UserRole
	LogLevel       pgx.LogLevel
	Logger         pgx.Logger
}

type databaseUsers struct {
	FullAccess *databaseUser `json:"fullAccess"`
}

type databaseUser struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// env ключ для каждого нужного значения для подключения ДБ
type databaseENV struct {
	Host           string
	Port           string
	DbName         string
	User           string
	Password       string
	SSLMode        string
	ConnectTimeout string
}

type databaseDSN struct {
	Name           string        `json:"name"`
	Host           string        `json:"host"`
	Port           int           `json:"port"`
	DbName         string        `json:"dbname"`
	SSLMode        string        `json:"sslmode"`
	ConnectTimeout int           `json:"connect_timeout"`
	User           databaseUsers `json:"user"`
}

type DSN struct {
	Name           string
	DBName         string
	Host           string
	Port           int
	ConnectTimeout int
	UserRole       int
	UserLogin      string
	UserPassword   string
	SSLMode        string
	// добавить сохранение сертификата
}

// unused
func ParseSSLMode(mode string) (SSLModeType, error) {
	modes := []string{
		string(SSLModeDisabled),
	}

	for _, m := range modes {
		if mode == m {
			return SSLModeType(mode), nil
		}
	}

	return "", fmt.Errorf(`wrong sslmode "%s", choose one from: %s`, mode, strings.Join(modes, ", "))
}
