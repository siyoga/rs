package logger

import (
	"container/heap"
	"encoding/json"
	"fmt"

	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

type driverZapLimitJSONSizeEncoder struct {
	zapcore.Encoder

	limit int
}

func newLimitJSONSizeEncoder(encoder zapcore.Encoder, limit int) *driverZapLimitJSONSizeEncoder {
	return &driverZapLimitJSONSizeEncoder{
		Encoder: encoder,
		limit:   limit,
	}
}

func (c *driverZapLimitJSONSizeEncoder) Clone() zapcore.Encoder {
	return &driverZapLimitJSONSizeEncoder{
		Encoder: c.Encoder.Clone(),
		limit:   c.limit,
	}
}

func (c *driverZapLimitJSONSizeEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	buf, err := c.Encoder.EncodeEntry(entry, fields)
	if err != nil {
		return buf, fmt.Errorf("encoder returned error: %w", err)
	}

	if buf.Len() < c.limit {
		return buf, nil
	}

	return c.reduceSize(buf), nil
}

func (c *driverZapLimitJSONSizeEncoder) reduceSize(buf *buffer.Buffer) *buffer.Buffer {
	const reduceRetries = 10

	content := map[string]interface{}{}

	if err := json.Unmarshal(buf.Bytes(), &content); err != nil {
		// can not unmarshal
		return buf
	}

	maxHeap := make(fieldMaxHeap, 0, len(content))

	removableLen := 0
	for k, v := range content {
		if systemFields[k] {
			continue
		}

		switch tv := v.(type) {
		case string:
			removableLen += len(tv)
			maxHeap = append(maxHeap, fieldLen{name: k, len: len(tv)})
		}
	}

	if len(maxHeap) == 0 {
		return buf
	}

	if buf.Len()-removableLen > c.limit {
		// even if we will remove all the fields, total size will still be greater then limit,
		// so just skip filtering as it is useless
		return buf
	}

	// O(n)
	heap.Init(&maxHeap)

	resultLen := buf.Len()
	for i := 0; i < reduceRetries; i++ {
		longest := maxHeap[0]

		if longest.len < 2 {
			// nothing to optimize, longest string is too short
			return buf
		}

		value := []rune(content[longest.name].(string))
		content[longest.name] = string(value[:len(value)/2])
		maxHeap[0].len = len(content[longest.name].(string))
		resultLen -= maxHeap[0].len

		if resultLen < c.limit {
			// estimated result is shorter than limit
			break
		}

		// O(log(n))
		heap.Fix(&maxHeap, 0)
	}

	result, err := json.Marshal(content)
	if err != nil {
		return buf
	}
	buf.Reset()
	_, _ = buf.Write(result)
	_ = buf.WriteByte('\n') // always end with newline

	return buf
}

type fieldLen struct {
	name string
	len  int
}

type fieldMaxHeap []fieldLen

func (f fieldMaxHeap) Len() int           { return len(f) }
func (f fieldMaxHeap) Less(i, j int) bool { return f[i].len > f[j].len } // this gives us Max Heap
func (f fieldMaxHeap) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }

func (f *fieldMaxHeap) Push(x interface{}) { *f = append(*f, x.(fieldLen)) }

func (f *fieldMaxHeap) Pop() interface{} {
	old := *f
	n := len(old)
	x := old[n-1]
	*f = old[0 : n-1]
	return x
}
