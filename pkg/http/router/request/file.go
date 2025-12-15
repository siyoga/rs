package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"

	goform "github.com/go-playground/form/v4"
)

const binaryContentType = "application/octet-stream"

type File struct {
	buf         io.Reader
	contentType string
	name        string
	size        int
}

type extensionCallback struct {
	callback   goform.DecodeCustomTypeFunc
	targetType interface{}
}

func NewFile(body []byte) File {
	return NewFileFromReader(bytes.NewReader(body), len(body))
}

func NewFileFromReader(reader io.Reader, size ...int) File {
	bodySize := 0
	if len(size) > 0 {
		bodySize = size[0]
	}

	return File{
		buf:         reader,
		contentType: binaryContentType,
		size:        bodySize,
	}
}

func (f *File) Read(p []byte) (n int, err error) {
	return f.buf.Read(p)
}

func (f *File) Name() string {
	return f.name
}

func (f *File) Size() int {
	return f.size
}

func (f *File) ContentType() string {
	return f.contentType
}

func (f *File) MarshalJSON() ([]byte, error) {
	body, err := io.ReadAll(f.buf)
	if err != nil {
		return nil, err
	}

	return json.Marshal(body)
}

func (f *File) UnmarshalJSON(val []byte) error {
	buf := []byte{}
	if err := json.Unmarshal(val, &buf); err != nil {
		return err
	}

	*f = NewFile(buf)

	return nil
}

func (*File) decodeUrlEncodedExtension() extensionCallback {
	return extensionCallback{
		callback: func(values []string) (interface{}, error) {
			return NewFile([]byte(values[0])), nil
		},
		targetType: File{},
	}
}

func (*File) decodeMultipartExtension(files map[string][]File) extensionCallback {
	return extensionCallback{
		callback: func(values []string) (interface{}, error) {
			val := values[0]
			delimIdx := strings.LastIndex(val, "|")
			if delimIdx == -1 {
				return nil, fmt.Errorf("there is no such file in file section")
			}

			key, idx := val[:delimIdx], val[delimIdx+1:]
			numIdx, err := strconv.Atoi(idx)
			if err != nil {
				return nil, fmt.Errorf("there is no such file in file section")
			}

			return files[key][numIdx], nil
		},
		targetType: File{},
	}
}
