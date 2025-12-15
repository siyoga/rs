package logger

import (
	"fmt"
	"go.uber.org/zap/zapcore"
	"time"
)

const (
	black color = iota + 30
	red
	green
	yellow
	blue
	magenta
	cyan
	white
	darkGray color = 90
)

var (
	debugLevel   = magenta.add("D")
	infoLevel    = blue.add("I")
	warningLevel = yellow.add("W")
	errorLevel   = red.add("E")
)

type color uint8

// Add adds the coloring to the given string.
func (c color) add(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", uint8(c), s)
}

type prettyEncoder struct {
	zapcore.Encoder
}

func newPrettyEncoder(cfg zapcore.EncoderConfig) zapcore.Encoder {
	cfg.EncodeTime = encodeTime
	cfg.EncodeLevel = encodeLevel
	cfg.ConsoleSeparator = " "
	return &prettyEncoder{zapcore.NewConsoleEncoder(cfg)}
}

func (p *prettyEncoder) AddString(key, value string) {
	if "env" == key || "tag" == key {
		return
	}
	p.Encoder.AddString(key, value)
}

func (p *prettyEncoder) Clone() zapcore.Encoder {
	return &prettyEncoder{p.Encoder.Clone()}
}

func encodeTime(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	const layout = time.TimeOnly

	enc.AppendString(darkGray.add(t.Format(layout)))
}

func encodeLevel(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(levelString(l))
}

func levelString(l zapcore.Level) string {
	switch l {
	case zapcore.DebugLevel:
		return debugLevel
	case zapcore.InfoLevel:
		return infoLevel
	case zapcore.WarnLevel:
		return warningLevel
	case zapcore.ErrorLevel:
		return errorLevel
	default:
		return red.add(l.CapitalString())
	}
}
