/*
 * MIT License
 *
 * Copyright (c) 2025 Peter Vrba
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package tester

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
)

// newRequest creates a new *request instance
func newRequest(method, path string, deps *Deps) *request {
	return &request{
		method: method,
		path:   path,
		deps:   deps,
		header: make(http.Header),
		query:  make(url.Values),
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
func (r *request) Body(t require.TestingT, body any) APIRequest {
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
func (r *request) Do(t require.TestingT, ctx context.Context) APIResponse {
	req, err := http.NewRequestWithContext(ctx, r.method, r.path, r.body)
	assert.NoErrorf(t, err, "failed to create request: %v", err)

	// check if we have headers
	if r.header == nil {
		r.header = make(http.Header)
	}

	// set content type to json
	r.header.Set(ContentTypeHeader, ContentTypeJSON)

	// set headers to request
	req.Header = r.header

	// prepare recorder for response
	rw := httptest.NewRecorder()

	// add query if present
	if len(r.query) > 0 {
		req.URL.RawQuery = r.query.Encode()
	}

	// add header if present
	if len(r.header) > 0 {
		req.Header = r.header
	}

	// do the request to the address or handler
	if r.deps.Address != "" {
		r.doAddress(t, rw, req)
	} else {
		r.doHandler(t, rw, req)
	}

	// prepare response
	result := newResponse(rw, req)

	// add body when present
	if rw.Body != nil {
		result.body = rw.Body.Bytes()
	}

	return result
}

// doAddress does the request to the address
func (r *request) doAddress(t require.TestingT, rw http.ResponseWriter, req *http.Request) {
	// prepare client without any timeout since we add context to the request
	var hc *http.Client
	if r.deps.Client != nil {
		hc = r.deps.Client
	} else {
		hc = &http.Client{}
	}

	// update host and scheme
	req.URL.Host = r.deps.Address
	req.URL.Scheme = "http"

	resp, err := hc.Do(req)
	require.NoErrorf(t, err, "failed to do request: %v", err)

	// write status code and headers
	rw.WriteHeader(resp.StatusCode)
	for key, values := range resp.Header {
		for _, value := range values {
			rw.Header().Add(key, value)
		}
	}

	// write body
	if resp.Body != nil {
		_, err := io.Copy(rw, resp.Body)
		assert.NoErrorf(t, err, "failed to copy body: %v", err)
	}
}

// doHandler does the request to the handler
func (r *request) doHandler(t require.TestingT, rw http.ResponseWriter, req *http.Request) {
	assert.NotPanicsf(t, func() {
		r.deps.Handler.ServeHTTP(rw, req)
	}, "handler panicked")
}

// Header sets the header of the request
func (r *request) Header(t require.TestingT, key, value string) APIRequest {
	if r.header == nil {
		r.header = make(http.Header)
	}
	r.header.Add(key, value)
	return r
}

// Query sets the query of the request
func (r *request) Query(t require.TestingT, key, value string) APIRequest {
	if r.query == nil {
		r.query = make(url.Values)
	}
	r.query.Add(key, value)
	return r
}
