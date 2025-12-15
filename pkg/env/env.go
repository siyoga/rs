package env

import (
	"fmt"
	"os"
	"strconv"
)

func LookupEnvString(key string) (string, error) {
	value, ok := os.LookupEnv(key)
	if !ok {
		return "", fmt.Errorf("environment variable %s not set", key)
	}

	return value, nil
}

func LookupEnvInt64(key string) (int64, error) {
	value, ok := os.LookupEnv(key)
	if !ok {
		return 0, fmt.Errorf("environment variable %s not set", key)
	}

	i, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("environment variable %s must be an integer", key)
	}

	return i, nil
}

func LookupEnvInt(key string) (int, error) {
	value, ok := os.LookupEnv(key)
	if !ok {
		return 0, fmt.Errorf("environment variable %s not set", key)
	}

	i, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("environment variable %s must be an integer", key)
	}

	return i, nil
}
