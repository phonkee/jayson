package tester

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func newRequest(method, path string, deps *Deps) *request {
	return &request{
		method: method,
		path:   path,
		deps:   deps,
	}
}

// request is the implementation of APIRequest
type request struct {
	method string
	path   string
	body   io.Reader
	deps   *Deps
	query  url.Values
	header http.Header
}

// Body sets the body of the request
func (r *request) Body(t *testing.T, body any) APIRequest {
	switch body.(type) {
	case string:
		r.body = strings.NewReader(body.(string))
	case []byte:
		r.body = strings.NewReader(string(body.([]byte)))
	case io.Reader:
		r.body = body.(io.Reader)
	default:
		buffer := new(bytes.Buffer)
		assert.NoErrorf(t, json.NewEncoder(buffer).Encode(body), "cannot marshal to json: %v", body)
		r.body = buffer
	}
	return r
}

// Do does the request and returns the response
func (r *request) Do(t *testing.T, ctx context.Context) APIResponse {
	req, err := http.NewRequestWithContext(ctx, r.method, r.path, r.body)
	assert.NoErrorf(t, err, "failed to create request: %v", err)

	// prepare recorder for response
	rw := httptest.NewRecorder()

	assert.NotPanicsf(t, func() {
		r.deps.Handler.ServeHTTP(rw, req)
	}, "handler panic")

	// add query if present
	if len(r.query) > 0 {
		req.URL.RawQuery = r.query.Encode()
	}

	// add header if present
	if len(r.header) > 0 {
		req.Header = r.header
	}

	// prepare response
	result := &response{
		rw:      rw,
		request: req,
	}

	// add body when present
	if rw.Body != nil {
		result.body = rw.Body.Bytes()
	}

	return result
}

// Header sets the header of the request
func (r *request) Header(t *testing.T, key, value string) APIRequest {
	if r.header == nil {
		r.header = make(http.Header)
	}
	r.header.Add(key, value)
	return r
}

// Query sets the query of the request
func (r *request) Query(t *testing.T, key, value string) APIRequest {
	if r.query == nil {
		r.query = make(url.Values)
	}
	r.query.Add(key, value)
	return r
}
