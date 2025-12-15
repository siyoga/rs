package request

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Request interface {
	URLPath() string
	Method() string
	JSONBody(data interface{}) error
}

type request struct {
	bodyContentType contentType

	urlPath string
	method  string
	query   url.Values
	headers http.Header
	cookies map[string]*http.Cookie

	body             io.ReadCloser
	multipartValues  map[string][]string
	multipartFiles   map[string][]File
	urlEncodedValues url.Values
}

var _ Request = (*request)(nil)

const contentTypeHeader = "content-type"

type contentType int

const (
	UndefinedDataBody contentType = iota
	MultiPartFormDataBody
	UrlEncodedFormBody
	JsonDataBody
)

func NewRequest(req *http.Request) (*request, error) {
	newRequest := &request{
		urlPath: req.URL.Path,
		method:  req.Method,
		headers: req.Header,
		body:    req.Body,

		cookies: make(map[string]*http.Cookie),
		query:   url.Values{},

		multipartValues:  make(map[string][]string),
		multipartFiles:   make(map[string][]File),
		urlEncodedValues: url.Values{},
	}

	newRequest.bodyContentType = newRequest.getContentTypeBody(req)
	newRequest.parseHttpRequest(req)

	if err := newRequest.parseMultipartForm(req); err != nil {
		return nil, err
	}

	if err := newRequest.parseUrlEncodedForm(req); err != nil {
		return nil, err
	}

	return newRequest, nil
}

func (r *request) URLPath() string {
	return r.urlPath
}

func (r *request) Method() string {
	return r.method
}

func (r *request) JSONBody(data interface{}) error {
	if r.body == nil || data == nil {
		return nil
	}

	if r.bodyContentType != JsonDataBody {
		return nil
	}

	decoder := json.NewDecoder(r.body)
	decoder.UseNumber()

	return decoder.Decode(data)
}

func (r *request) getContentTypeBody(req *http.Request) contentType {
	contentTypeHeaderValue := req.Header.Get(contentTypeHeader)
	if strings.Contains(contentTypeHeaderValue, "multipart/form-data") {
		return MultiPartFormDataBody
	} else if strings.Contains(contentTypeHeaderValue, "application/x-www-form-urlencoded") {
		return UrlEncodedFormBody
	} else if strings.Contains(contentTypeHeaderValue, "application/json") {
		return JsonDataBody
	} else {
		return UndefinedDataBody
	}
}

func (r *request) parseHttpRequest(req *http.Request) {
	for _, cookie := range req.Cookies() {
		r.cookies[cookie.Name] = cookie
	}

	r.query = normalizeValues(req.URL.Query())
	for queryKey := range r.query {
		if len(queryKey) == 0 {
			r.query.Del(queryKey)
		}
	}
}

func normalizeValues(values url.Values) url.Values {
	normalizedValues := make(url.Values)

	for key, value := range values {
		newKey := key

		// Для поддержки array[]=1&array[]=2 убираем [] в конце ключа
		n := len(newKey)
		if n > 2 && newKey[n-2] == '[' && newKey[n-1] == ']' {
			newKey = newKey[:n-2]
		}

		for i := range value {
			normalizedValues.Add(newKey, value[i])
		}
	}

	return normalizedValues
}
