package logger

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"runtime"
	"strings"
	"sync"
	"time"
)

type driver interface {
	Debug(ctx context.Context, arg ...interface{})
	Info(ctx context.Context, arg ...interface{})
	Warning(ctx context.Context, arg ...interface{})
	Error(ctx context.Context, arg ...interface{})
	Flush(timeout time.Duration) error
}

type zapDriver struct {
	*options

	zap *zap.Logger
}

func newZapDriver(o *options) (*zapDriver, error) {
	level, err := parseZapLevel(o.level)
	if err != nil {
		return nil, fmt.Errorf("can't parse log level %q: %w", o.level, err)
	}

	var encoder zapcore.Encoder

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       keyOrOmit(fieldTime, o.includedFields),
		LevelKey:      keyOrOmit(fieldLevel, o.includedFields),
		NameKey:       keyOrOmit(fieldLogger, o.includedFields),
		CallerKey:     keyOrOmit(fieldCaller, o.includedFields),
		FunctionKey:   zapcore.OmitKey,
		MessageKey:    keyOrOmit(fieldMessage, o.includedFields),
		StacktraceKey: keyOrOmit(fieldStacktrace, o.includedFields),
		EncodeLevel: func(level zapcore.Level, encoder zapcore.PrimitiveArrayEncoder) {
			encoder.AppendString(string(level.CapitalString()[0]))
		},
		EncodeTime:     zapcore.RFC3339TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	format := strings.ToLower(o.format)

	switch format {
	case FormatJSON:
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	case FormatText:
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	case FormatPretty:
		encoder = newPrettyEncoder(encoderConfig)
	}

	if o.maxLogEntrySize > 0 && format == FormatJSON {
		encoder = newLimitJSONSizeEncoder(encoder, o.maxLogEntrySize)
	}

	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(o.output),
		level,
	)

	if len(o.includedFields) > 0 {
		core = newDriverZapFilteringCore(core, o.includedFields)
	}

	fields := make([]interface{}, 0, len(o.defaultFields)+4)
	for k, v := range o.defaultFields {
		fields = append(fields, k, v)
	}

	zapLogger := zap.New(core).Sugar().With(fields...).Desugar()

	return &zapDriver{
		options: o,
		zap:     zapLogger,
	}, nil
}

func parseZapLevel(level string) (zapcore.Level, error) {
	upperLevel := strings.ToUpper(level)

	switch upperLevel {
	case "DEBUG", "D":
		return zapcore.DebugLevel, nil
	case "INFO", "I":
		return zapcore.InfoLevel, nil
	case "WARN", "W", "WARNING":
		return zapcore.WarnLevel, nil
	case "ERROR", "E", "ERR", "F", "FATAL", "P", "PANIC":
		return zapcore.ErrorLevel, nil
	default:
		return zapcore.InvalidLevel, fmt.Errorf("unknown level: %s", level)
	}
}

func keyOrOmit(key string, included []string) string {
	if len(included) == 0 {
		return key
	}

	for i := range included {
		if included[i] == key {
			return key
		}
	}

	return zapcore.OmitKey
}

func (d *zapDriver) Debug(ctx context.Context, arg ...interface{}) {
	d.logger(ctx).Debug(arg...)
}

func (d *zapDriver) Info(ctx context.Context, arg ...interface{}) {
	d.logger(ctx).Info(arg...)
}

func (d *zapDriver) Warning(ctx context.Context, arg ...interface{}) {
	d.logger(ctx).Warn(arg...)
}

func (d *zapDriver) Error(ctx context.Context, arg ...interface{}) {
	d.logger(ctx).Error(arg...)
}

func (d *zapDriver) Flush(timeout time.Duration) error {
	return d.zap.Sync()
}

const (
	defSkipDepth = 3
	maxSkipDepth = 12
)

var loggerFieldsPool = sync.Pool{New: func() interface{} { return &[]zap.Field{} }}

func (d *zapDriver) logger(ctx context.Context) *zap.SugaredLogger {
	loggerFieldsPtr := loggerFieldsPool.Get().(*[]zap.Field)
	loggerFields := *loggerFieldsPtr

	defer func() {
		loggerFields = loggerFields[:0]
		*loggerFieldsPtr = loggerFields
		loggerFieldsPool.Put(loggerFieldsPtr)
	}()

	if callFields, ok := getCallInfo(ctx); ok {
		loggerFields = append(loggerFields, zap.Int(callFields[0].(string), callFields[1].(int)))
		loggerFields = append(loggerFields, zap.String(callFields[2].(string), callFields[3].(string)))
	}

	fields, put := fieldsPooled(ctx)
	for k, v := range fields {
		loggerFields = append(loggerFields, zap.Any(k, v))
	}
	put()

	tags, put := tagsPooled(ctx)
	for k, v := range tags {
		loggerFields = append(loggerFields, zap.Any(k, v))
	}
	put()

	return d.zap.With(loggerFields...).Sugar()
}

// getCallInfo returns fields with execution context like functions and line number
func getCallInfo(ctx context.Context) ([4]interface{}, bool) {
	// if pc is passed withing context, then just use it
	if pc, ok := pcField(ctx); ok && pc > 0 {
		fs := runtime.CallersFrames([]uintptr{pc})
		f, _ := fs.Next()
		return [...]interface{}{fieldLineNo, f.Line, fieldFunction, f.Function}, true
	}

	var (
		pc   uintptr
		line int
		ok   bool
		f    string
	)

	// Skip all calls that should not be included into stack trace
	// This helps us to see actual place where call was made, but not internal logger place
	for i := defSkipDepth; i < maxSkipDepth; i++ {
		pc, _, line, ok = runtime.Caller(i)
		if !ok {
			return [...]interface{}{nil, nil, nil, nil}, false
		}
		f = runtime.FuncForPC(pc).Name()
	}
	return [...]interface{}{fieldLineNo, line, fieldFunction, f}, true
}
