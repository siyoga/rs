package logger

import "go.uber.org/zap/zapcore"

type driverZapFilteringCore struct {
	includedFields map[string]bool

	inner zapcore.Core
}

func newDriverZapFilteringCore(inner zapcore.Core, includedFields []string) *driverZapFilteringCore {
	core := &driverZapFilteringCore{
		inner:          inner,
		includedFields: make(map[string]bool, len(includedFields)),
	}

	for i := range includedFields {
		core.includedFields[includedFields[i]] = true
	}

	return core
}

func (z *driverZapFilteringCore) Enabled(level zapcore.Level) bool {
	return z.inner.Enabled(level)
}

func (z *driverZapFilteringCore) With(fields []zapcore.Field) zapcore.Core {
	inner := z.inner.With(z.filterFields(fields))

	return &driverZapFilteringCore{
		includedFields: z.includedFields,
		inner:          inner,
	}
}

func (z *driverZapFilteringCore) Check(entry zapcore.Entry, checked *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if innerChecked := z.inner.Check(entry, checked); innerChecked != nil {
		// if inner decided to log this entity, then let's log it but with our core
		return checked.AddCore(entry, z)
	}
	return nil
}

func (z *driverZapFilteringCore) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	return z.inner.Write(entry, z.filterFields(fields))
}

func (z *driverZapFilteringCore) Sync() error {
	return z.inner.Sync()
}

func (z *driverZapFilteringCore) filterFields(fields []zapcore.Field) []zapcore.Field {
	pos := 0
	for i := range fields {
		if len(z.includedFields) > 0 && !z.includedFields[fields[i].Key] {
			continue // field is not in the list of enabled ones
		}

		fields[pos] = fields[i]
		pos++
	}
	return fields[:pos]
}
