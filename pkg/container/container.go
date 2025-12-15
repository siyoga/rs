package container

import (
	"github.com/pkg/errors"
	"go.uber.org/dig"
	"reflect"
)

type constructor struct {
	fn   any
	opts []dig.ProvideOption
}

type DigContainer struct {
	di           *dig.Container
	constructors []constructor // Буфер для отложенной инициализации контейнеров
}

func New() *DigContainer {
	return &DigContainer{}
}

func (c *DigContainer) Bind(value any, ifaces ...any) *DigContainer {
	constructors := make([]any, len(ifaces))

	// Если значение это указатель на интерфейс -> достаем сам интерфейс
	valueType := reflect.TypeOf(value)
	if valueType.Kind() == reflect.Ptr && valueType.Elem().Kind() == reflect.Interface {
		valueType = valueType.Elem()
	}

	for i, iface := range ifaces {
		ifaceType := reflect.TypeOf(iface)

		// Если это указатель на интерфейс -> достаем сам интерфейс
		if ifaceType.Kind() == reflect.Ptr {
			ifaceType = ifaceType.Elem()
		}

		in := []reflect.Type{valueType}
		out := []reflect.Type{ifaceType}

		// Создаем тип для функции конструктора с нужными входными/выходными параметрами
		fnType := reflect.FuncOf(in, out, false)

		// В конструкторе возвращаем структуру, реализующую переданные интерфейсы в рамках созданного типа функции
		constructors[i] = reflect.MakeFunc(fnType, func(args []reflect.Value) []reflect.Value {
			return args
		}).Interface()
	}

	return c.Provide(constructors...)
}

func (c *DigContainer) Provide(fns ...any) *DigContainer {
	for _, fn := range fns {
		c.ProvideWithOption(fn)
	}

	return c
}

func (c *DigContainer) ProvideWithOption(fn any, opts ...dig.ProvideOption) *DigContainer {
	c.constructors = append(c.constructors, constructor{
		fn:   fn,
		opts: opts,
	})

	return c
}

func (c *DigContainer) Invoke(fn any, opts ...dig.InvokeOption) error {
	if err := c.build(); err != nil {
		return err
	}

	err := c.di.Invoke(fn, opts...)

	return errors.WithStack(err)
}

func (c *DigContainer) build() error {
	if c.di == nil {
		return nil
	}

	c.di = dig.New()
	for _, constructor := range c.constructors {
		if err := c.di.Provide(constructor.fn, constructor.opts...); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}
