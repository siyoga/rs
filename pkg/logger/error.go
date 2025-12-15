package logger

import (
	"context"
	"errors"
)

type errorWithFields interface {
	LoggerFields() map[string]any
}

type errWrapper struct {
	fields map[string]interface{}

	err error
}

func (e *errWrapper) Error() string {
	return e.err.Error()
}

func (e *errWrapper) Unwrap() error {
	return e.err
}

func (e *errWrapper) Cause() error {
	return e.err
}

func withError(ctx context.Context, err error) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	if err == nil {
		return ctx
	}

	var (
		errFields = map[string]interface{}{}
	)

	var wrappedErr *errWrapper
	if errors.As(err, &wrappedErr) {
		// if we have errWrapper somewhere inside err, then we will extract its fields and tags
		errFields = wrappedErr.fields
	}

	if wrappedErr, ok := err.(*errWrapper); ok {
		// if this err is already err wrapper, then we will unwrap it to reduce stack.
		// Fields and tags were extracted in the previous step
		err = wrappedErr.err
	}

	if errWithFields, ok := err.(errorWithFields); ok {
		for k, v := range errWithFields.LoggerFields() {
			errFields[k] = v
		}
	}

	ctx = withFields(ctx, errFields)

	return withField(ctx, errorValueKey, err)
}

// wrapError оборачивает переданную ошибку err тегами и полями из ctx и возвращает новую ошибку,
// которую затем можно использовать в методах withField и подобных для логирования ее вместе с данными из контекста
func wrapError(ctx context.Context, err error) error {
	if err == nil {
		return err // maintain error type
	}

	if ctx == nil {
		return err
	}

	var (
		ctxFields = fieldsSlow(ctx)
	)

	var wrappedErr *errWrapper
	if errors.As(err, &wrappedErr) {
		// if we have errWrapper somewhere inside err, then we will extract its fields and tags
		for name, value := range wrappedErr.fields {
			if _, ok := ctxFields[name]; !ok {
				ctxFields[name] = value
			}
		}
	}

	if wrappedErr, ok := err.(*errWrapper); ok {
		// if this err is already err wrapper, then we will unwrap it to reduce stack.
		// Fields and tags were extracted in the previous step
		err = wrappedErr.err
	}

	return &errWrapper{
		fields: ctxFields,
		err:    err,
	}
}
