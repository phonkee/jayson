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
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
)

// newRequest creates a new *request instance
func newRequest(t TestingT, method, path string, deps *Deps) *request {
	req, err := http.NewRequest(method, path, nil)
	require.NoError(t, err, "failed to create request: %v", err)

	// set content type to json
	req.Header.Set(ContentTypeHeader, ContentTypeJSON)

	return &request{
		deps: deps,
		req:  req,
	}
}

// request is the implementation of APIRequest
type request struct {
	deps *Deps
	req  *http.Request
}

// Body sets the body of the request
func (r *request) Body(t require.TestingT, body any) APIRequest {
	switch body.(type) {
	case string:
		r.req.Body = NewReadCloser(strings.NewReader(body.(string)))
	case []byte:
		r.req.Body = NewReadCloser(strings.NewReader(string(body.([]byte))))
	case io.Reader:
		r.req.Body = NewReadCloser(body.(io.Reader))
	default:
		if body == nil {
			r.req.Body = nil
		} else {
			buffer := new(bytes.Buffer)
			assert.NoErrorf(t, json.NewEncoder(buffer).Encode(body), "cannot marshal to json: %v", body)
			r.req.Body = NewReadCloser(buffer)
		}
	}
	return r
}

// Do the request and returns the response
func (r *request) Do(t require.TestingT, ctx context.Context) APIResponse {
	// add context to request
	req := r.req.WithContext(ctx)

	// prepare recorder for response
	rw := httptest.NewRecorder()

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
func (r *request) Header(_ require.TestingT, key, value string) APIRequest {
	r.req.Header.Add(key, value)
	return r
}

// Print prints the response to given writer, if not given, it prints to stdout
func (r *request) Print(writer ...io.Writer) APIRequest {
	var w io.Writer
	if len(writer) > 0 && writer[0] != nil {
		w = writer[0]
	} else {
		w = os.Stdout
	}

	// printf is a helper function to print formatted output and ignore errors
	printf := func(format string, args ...any) {
		_, _ = fmt.Fprintf(w, format+"\n", args...)
	}
	// TODO: do this in a better way
	printf("Print:")
	printf("  Request:")
	printf("    URL: %v", r.req.URL.String())
	printf("    Method: %v", r.req.Method)
	printf("    Header: %v", r.req.Header)

	return r

}

// Query sets the query of the request
func (r *request) Query(t require.TestingT, key, value string) APIRequest {
	q := r.req.URL.Query()
	q.Add(key, value)
	r.req.URL.RawQuery = q.Encode()
	return r
}
