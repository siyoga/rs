package logger

import "io"

type options struct {
	output io.Writer

	includedFields []string
	defaultFields  map[string]interface{}

	app             string
	level           string
	enabled         bool
	format          string
	maxLogEntrySize int
}

type Option func(*options)

func newOptions(opts ...Option) *options {
	opt := &options{
		output: envGetOutput(),

		defaultFields:  make(map[string]interface{}),
		includedFields: envStringArray(envIncludeFields, nil),

		app:     envString(envAppName, ""),
		level:   envString(envLevel, defaultLevel),
		enabled: envBool(envEnabled, true),
		format:  envString(envOutput, defaultFormat),
	}

	for _, o := range opts {
		o(opt)
	}

	return opt
}
