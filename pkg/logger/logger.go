package logger

import (
	"context"
	"fmt"
	"time"
)

type Logger struct {
	*options
	driver driver
}

func New(opts ...Option) (*Logger, error) {
	o := newOptions(opts...)

	d, err := createDriver(o)
	if err != nil {
		return nil, err
	}

	return &Logger{
		options: o,
		driver:  d,
	}, nil

}

func createDriver(o *options) (driver, error) {
	zapDriver, err := newZapDriver(o)
	if err != nil {
		return nil, fmt.Errorf("create zap driver: %w", err)
	}

	return zapDriver, nil
}

// Debug логирует сообщение с уровнем debug
func (l *Logger) Debug(ctx context.Context, args ...interface{}) {
	ctx, args = withArgs(ctx, args...)
	l.driver.Debug(ctx, args...)
}

// Info логирует сообщение с уровнем info
func (l *Logger) Info(ctx context.Context, args ...interface{}) {
	ctx, args = withArgs(ctx, args...)
	l.driver.Info(ctx, args...)
}

// Warning логирует сообщение с уровнем warning
func (l *Logger) Warning(ctx context.Context, args ...interface{}) {
	ctx, args = withArgs(ctx, args...)
	l.driver.Warning(ctx, args...)
}

// Error логирует сообщение с уровнем error и отправляет ошибку в Sentry
func (l *Logger) Error(ctx context.Context, args ...interface{}) {
	ctx, args = withArgs(ctx, args...)
	l.driver.Error(ctx, args...)
}

// WithField добавляет в контекст поле ключ/значение. При передаче возвращаемого
// контекста в один из методов логирования данное поле будет добавлено в лог.
func (l *Logger) WithField(ctx context.Context, k string, v interface{}) context.Context {
	return withField(ctx, k, v)
}

// WithFields добавляет в контекст набор полей ключ/значение. При передаче возвращаемого
// контекста в один из методов логирования данное поле будет добавлено в лог.
func (l *Logger) WithFields(ctx context.Context, fields map[string]interface{}) context.Context {
	return withFields(ctx, fields)
}

// WithError - обёртка над методом WithField, добавляет в контекст указанную ошибку,
// которая будет выведена в лог по ключу `error`
// Аналогом данного метода является вызов WithField(ctx, "error", err)
//
// Если передать в метод в качестве ошибки результат вызова WrapError(ctx, err), то помимо
// ошибки также будет залогированы поля и теги из исходного контекста
//
// Если ошибка реализует метод LoggerFields() map[string]any, то он также будет вызван для получения полей.
// Если ошибка реализует метод LoggerTags() map[string]string, то он также будет вызван для получения тегов.
func (l *Logger) WithError(ctx context.Context, err error) context.Context {
	return withError(ctx, err)
}

// Fields возвращает все добавленные через методы WithField, WithFields и WithError в context.Context поля
func (l *Logger) Fields(ctx context.Context) map[string]interface{} {
	return fieldsSlow(ctx)
}

// WithContext копируют все добавленные через методы WithFields
// теги и поля из src в новый контекст на основе ctx
func (l *Logger) WithContext(ctx context.Context, src context.Context) context.Context {
	return withContext(ctx, src)
}

// WrapError оборачивает переданную ошибку err тегами и полями из ctx и возвращает новую ошибку,
// которую затем можно использовать в методе WithError для логирования ее вместе с данными из контекста.
//
// Если в контексте и в ошибке содержатся одинаковые поля и/или теги, то будут использованы значения
// из контекста.
func (l *Logger) WrapError(ctx context.Context, err error) error {
	return wrapError(ctx, err)
}

// Field позволяет передать в логгер поле ключ/значение.
func (l *Logger) Field(k string, v any) any {
	return field{k: k, v: &v}
}

// FieldErr позволяет передать в логгер поле с типом error.
func (l *Logger) FieldErr(err error) any {
	return field{err: err}
}

// Tag позволяет передать в логгер тег ключ/значение.
func (l *Logger) Tag(k string, v string) any {
	return field{k: k, tv: &v}
}

// Flush отправит оставшиеся в буфере логи в sentry в течение переданного таймаута
// вернёт false в случае, если не удалось отправить все логи в течении этого таймаута
func (l *Logger) Flush(timeout time.Duration) bool {
	if err := l.driver.Flush(timeout); err != nil {
		return false
	}

	return true
}
