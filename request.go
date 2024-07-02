package scraper

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type HttpMethod string

const (
	GET    HttpMethod = "GET"
	POST   HttpMethod = "POST"
	PUT    HttpMethod = "PUT"
	PATCH  HttpMethod = "PATCH"
	DELETE HttpMethod = "DELETE"
)

type FormData map[string]string

type Request interface {
	WithData(data FormData) Request
	WithContext(ctx context.Context) Request
	WithJson(jsonStr string) Request
	WithHeaders(headers map[string]string) Request
	WithBearerToken(token string) Request
	Do() (Response, error)
}

type request struct {
	targetUrl string
	method    HttpMethod

	data     FormData
	withData bool

	jsonStr  string
	withJson bool

	ctx         context.Context
	withContext bool

	headers     map[string]string
	withHeaders bool

	request *http.Request
}

var _ Request = (*request)(nil)

func WithData(data FormData) Option {
	return func(r *request) {
		r.data = data
		r.withData = true
		formData := url.Values{}
		for k, v := range data {
			formData.Set(k, v)
		}
		r.request.Body = io.NopCloser(strings.NewReader(formData.Encode()))
		r.request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
}

func (r *request) WithData(data FormData) Request {
	WithData(data)(r)
	return r
}

func WithJson(jsonStr string) Option {
	return func(r *request) {
		r.jsonStr = jsonStr
		r.withJson = true
		r.request.Header.Add("Content-Type", "application/json")
		r.request.Body = io.NopCloser(bytes.NewBuffer([]byte(jsonStr)))
	}
}

func (r *request) WithJson(jsonStr string) Request {
	WithJson(jsonStr)(r)
	return r
}

func WithHeaders(headers map[string]string) Option {
	return func(r *request) {
		r.AddHeaders(headers)
	}
}

func (r *request) WithHeaders(headers map[string]string) Request {
	WithHeaders(headers)(r)
	return r
}

func WithContext(ctx context.Context) Option {
	return func(r *request) {
		r.ctx = ctx
		r.withContext = true
		r.request = r.request.WithContext(ctx)
	}
}

func (r *request) WithContext(ctx context.Context) Request {
	WithContext(ctx)(r)
	return r
}

func (r *request) AddHeader(k, v string) {
	if r.headers == nil {
		r.headers = make(map[string]string)
	}
	r.headers[k] = v
	r.withHeaders = true
}

func (r *request) AddHeaders(hs map[string]string) {
	if r.headers == nil {
		r.headers = hs
	} else {
		for k, v := range hs {
			r.headers[k] = v
		}
	}
	r.withHeaders = true
}

func WithBearerToken(token string) Option {
	return func(r *request) {
		r.AddHeader("Authorization", fmt.Sprintf("Bearer %s", token))
	}
}

func (r *request) WithBearerToken(token string) Request {
	WithBearerToken(token)(r)
	return r
}

func (r *request) Do() (Response, error) {
	var err error
	var httpResponse *http.Response

	// Apply all headers
	if r.withHeaders {
		for k, v := range r.headers {
			r.request.Header.Add(k, v)
		}
	}
	httpResponse, err = http.DefaultClient.Do(r.request)
	resp := NewResponse(httpResponse)

	if err != nil {
		err = fmt.Errorf("failed to do http request: %s", err.Error())
		log.Println(err.Error())
		return nil, err
	}

	return resp, nil
}

type Option func(r *request)

func NewRequest(targetUrl string, method HttpMethod, options ...Option) (Request, error) {
	var err error
	var httpRequest *http.Request
	httpRequest, err = http.NewRequest(string(method), targetUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %s", err.Error())
	}
	r := &request{
		method:    method,
		targetUrl: targetUrl,
		request:   httpRequest,
	}

	for _, opt := range options {
		opt(r)
	}

	return r, nil
}
