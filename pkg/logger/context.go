package logger

import (
	"context"
	"sync"
)

type key int

const (
	fieldsKey key = iota
	pcKey

	errorValueKey = "error"

	maxValueDepth = 2000
)

type value struct {
	prev *value
	key  string
	val  interface{}
	// level is a counter that defines current size of value list.
	// it could be zero event for the first element, because compactValue method
	// use it to stop optimization when it reaches zero
	level int
}

type field struct {
	k   string
	v   *any
	tv  *string
	err error
}

func pcField(ctx context.Context) (uintptr, bool) {
	if ctx == nil {
		return 0, false
	}

	pc, ok := ctx.Value(pcKey).(uintptr)
	if !ok {
		return 0, false
	}

	return pc, true
}

func withArgs(ctx context.Context, args ...any) (context.Context, []any) {
	if ctx == nil {
		ctx = context.Background()
	}

	if args == nil {
		return ctx, args
	}

	var (
		// do not allocate filtered in advance: if where are no filed args, then this is useless
		filtered []any
		pos      int
	)
	for i := range args {
		if f, ok := args[i].(field); ok {
			if filtered == nil { // allocate filtered only here
				filtered = make([]any, len(args))
				copy(filtered, args) // copy all content
				pos = i              // mark last simple arg position
			}
			ctx = withInnerField(ctx, f)
			continue
		}

		if filtered != nil {
			filtered[pos] = args[i]
			pos++
		}
	}

	if filtered != nil {
		return ctx, filtered[:pos]
	}

	return ctx, args
}

func withInnerField(ctx context.Context, f field) context.Context {
	switch {
	case f.v != nil:
		return withField(ctx, f.k, *f.v)
	case f.err != nil:
		return withError(ctx, f.err)
	default:
		return ctx
	}
}

func withField(ctx context.Context, k string, v interface{}) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	next := appendValue(&value{
		key: k,
		val: v,
	}, ctx.Value(fieldsKey), maxValueDepth)

	return context.WithValue(ctx, fieldsKey, next)
}

func withFields(ctx context.Context, fields map[string]interface{}) context.Context {
	for key, value := range fields {
		ctx = withField(ctx, key, value)
	}
	return ctx
}

func fieldsSlow(ctx context.Context) map[string]interface{} {
	res := map[string]interface{}{}
	fieldsToMap(ctx, &res)
	return res
}

func appendValue(v *value, maybePrev interface{}, maxDepth int) *value {
	v.level = 0
	if prev, ok := maybePrev.(*value); ok && prev != nil {
		v.prev = prev
		v.level = prev.level + 1
	}

	compactValue(v, maxDepth)

	return v
}

// compactValue пытается оптимизировать размер структуры value.
func compactValue(v *value, maxDepth int) {
	if v.level < maxDepth {
		return
	}

	// rebuild
	values := make([]value, v.level)
	head := &value{}
	tail := head
	found := make(map[string]bool, maxDepth/10) // estimate that there will be 90% of the same keys -> 10% keys are unique
	pos := 0
	for cur := v; cur != nil; cur = cur.prev {
		if found[cur.key] {
			continue
		}
		found[cur.key] = true

		var newVal *value
		if pos < len(values) {
			newVal = &values[pos]
		} else {
			newVal = &value{}
		}
		pos++

		// every newVal receives newVal.level equal to zero.
		// this will indicate that optimization has been already performed
		newVal.key = cur.key
		newVal.val = cur.val

		tail.prev = newVal
		tail = newVal

		if cur.level == 0 {
			tail.prev = cur.prev
			break
		}
	}

	*v = *head.prev
}

var fieldsPool = sync.Pool{New: func() interface{} { return &map[string]interface{}{} }}

func fieldsPooled(ctx context.Context) (map[string]interface{}, func()) {
	resPtr := fieldsPool.Get().(*map[string]interface{})
	put := func() {
		clear(*resPtr)
		fieldsPool.Put(resPtr)
	}

	fieldsToMap(ctx, resPtr)

	return *resPtr, put
}

func fieldsToMap(ctx context.Context, dst *map[string]interface{}) {
	if ctx == nil {
		return
	}

	v, ok := ctx.Value(fieldsKey).(*value)
	if !ok {
		return
	}

	for cur := v; cur != nil; cur = cur.prev {
		if _, ok := (*dst)[cur.key]; !ok {
			(*dst)[cur.key] = cur.val
		}
	}
}

func withContext(ctx context.Context, src context.Context) context.Context {
	flds, put := fieldsPooled(src)
	for k, v := range flds {
		ctx = withField(ctx, k, v)
	}
	put()
	return ctx
}
