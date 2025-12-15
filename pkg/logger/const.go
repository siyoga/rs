package logger

import (
	"os"
	"strconv"
	"strings"
)

const pkgLogger = "github.com/siyoga/rollstory/logger"

// ENV
const (
	envAppName       = "APP_NAME"
	envEnvironment   = "ENVIRONMENT"
	envEnabled       = "LOGGER_ENABLED"
	envOutput        = "LOGGER_OUTPUT"
	envLevel         = "LOGGER_LEVEL"
	envIncludeFields = "LOGGER_INCLUDE_FIELDS"
	envExcludeFields = "LOGGER_EXCLUDE_FIELDS"
)

// LOGGER FIELDS
const (
	fieldTime       = "time"
	fieldLevel      = "level"
	fieldType       = "type"
	fieldLogger     = "logger"
	fieldMessage    = "message"
	fieldCaller     = "caller"
	fieldStacktrace = "stacktrace"
	fieldTracing    = "_tracing"
	fieldLineNo     = "lineno"
	fieldFunction   = "function"
	fieldUserID     = "user_id"
)

// DEFAULTS
const (
	defaultLevel  = "INFO"
	defaultFormat = "JSON"
	defaultOutput = "STDERR"
)

// FORMATS
const (
	FormatJSON   = "json"
	FormatText   = "text"
	FormatPretty = "pretty"
)

// SYS FIELDS
var (
	systemFields = map[string]bool{
		fieldTime:       true,
		fieldType:       true,
		fieldLogger:     true,
		fieldCaller:     true,
		fieldMessage:    true,
		fieldStacktrace: true,
		fieldUserID:     true,
	}
	defaultIgnoredPkgs = []string{"pkgLogger"}
)

func envString(env string, def string) string {
	if s := os.Getenv(env); s != "" {
		return s
	}
	return def
}

func envStringArray(env string, def []string) []string {
	val := strings.TrimSpace(os.Getenv(env))
	if val == "" {
		return def
	}

	res := strings.Split(val, ",")
	for i := range res {
		res[i] = strings.TrimSpace(res[i])
	}
	return res
}

func envBool(env string, def bool) bool {
	b, err := strconv.ParseBool(os.Getenv(env))
	if err != nil {
		return def
	}
	return b
}

func envGetOutput() *os.File {
	outputKey := strings.ToLower(envString(envOutput, defaultOutput))
	switch outputKey {
	case "stdout":
		return os.Stdout
	case "stderr":
		fallthrough
	default:
		return os.Stderr
	}
}
