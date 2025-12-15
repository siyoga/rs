package response

import (
	"encoding/json"
	"io"
	"net/http"
)

type response struct {
	Headers       map[string]string
	Cookies       map[string]string
	NativeCookies []http.Cookie
	Code          int
	RawBody       io.Reader
	JsonBody      interface{}
}

type Response interface {
	Builder
	Sender
}

type Builder interface {
	SetHeader(key string, value string) Response
	SetCookie(key, value string) Response
	SetNativeCookie(cookie http.Cookie) Response
	SetStatusCode(code int) Response
	SetRawBody(body io.Reader) Response
	SetJsonBody(body interface{}) Response
}

type Sender interface {
	Send(http.ResponseWriter) error
}

func NewResponse() Response {
	return &response{
		Headers:       make(map[string]string),
		Cookies:       make(map[string]string),
		NativeCookies: make([]http.Cookie, 0),
		Code:          http.StatusOK,
	}
}

func (r *response) SetHeader(key string, value string) Response {
	r.Headers[key] = value
	return r
}

func (r *response) SetCookie(key, value string) Response {
	r.Cookies[key] = value
	return r
}

func (r *response) SetNativeCookie(cookie http.Cookie) Response {
	r.NativeCookies = append(r.NativeCookies, cookie)
	return r
}

func (r *response) SetStatusCode(code int) Response {
	r.Code = code
	return r
}

func (r *response) SetRawBody(body io.Reader) Response {
	r.RawBody = body
	return r
}

func (r *response) SetJsonBody(body interface{}) Response {
	r.JsonBody = body
	return r
}

func (r *response) Send(w http.ResponseWriter) error {
	if len(r.Headers) > 0 {
		for headerName, headerValue := range r.Headers {
			w.Header().Set(headerName, headerValue)
		}
	}

	if len(r.Cookies) > 0 {
		for key, value := range r.Cookies {
			http.SetCookie(w, &http.Cookie{
				Name:  key,
				Value: value,
			})
		}
	}
	if len(r.NativeCookies) > 0 {
		for i := range r.NativeCookies {
			http.SetCookie(w, &r.NativeCookies[i])
		}
	}

	w.WriteHeader(r.Code)

	if r.JsonBody != nil {
		w.Header().Set("Content-Type", "application/json")
		return json.NewEncoder(w).Encode(r.JsonBody)
	}

	if r.RawBody != nil {
		_, err := io.Copy(w, r.RawBody)
		return err
	}

	return nil
}
