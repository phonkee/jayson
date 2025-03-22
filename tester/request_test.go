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
	"context"
	"github.com/phonkee/jayson/tester/action"
	"github.com/phonkee/jayson/tester/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"net/http"
	"strings"
	"testing"
)

// must panics if error is not nil
func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

// noopHandler does nothing for testing purposes
type noopHandler struct {
	status int
}

// ServeHTTP is handler that writes status code
func (n noopHandler) ServeHTTP(rw http.ResponseWriter, _ *http.Request) {
	if n.status == 0 {
		rw.WriteHeader(http.StatusOK)
	} else {
		rw.WriteHeader(n.status)
	}
}

func TestRequest_Header(t *testing.T) {
	r := newRequest(t, http.MethodGet, "http://localhost:8080", &Deps{Handler: noopHandler{}}).
		Header(t, "key", "action")
	r.Do(t, context.Background())
	assert.Equal(t, []string{ContentTypeJSON}, r.(*request).req.Header.Values(ContentTypeHeader))
	assert.Equal(t, ContentTypeJSON, r.(*request).req.Header.Get(ContentTypeHeader))
	assert.Equal(t, "action", r.(*request).req.Header.Get("key"))
}

func TestRequest_Query(t *testing.T) {
	for _, item := range []struct {
		name     string
		url      string
		query    map[string]string
		expected string
	}{
		{
			name:     "test query",
			query:    map[string]string{"key": "action"},
			expected: "key=action",
		},
	} {
		t.Run(item.name, func(t *testing.T) {
			rt := mocks.NewRoundTripper(t)

			rt.On("RoundTrip", mock.MatchedBy(func(req *http.Request) bool {
				return req.URL.RawQuery == item.expected
			})).Return(&http.Response{
				StatusCode: http.StatusOK,
			}, nil)

			r := newRequest(
				t,
				http.MethodGet,
				item.url,
				&Deps{
					Address: "localhost:8080",
					Client: &http.Client{
						Transport: rt,
					},
				},
			)
			for key, value := range item.query {
				r.Query(t, key, value)
			}
			r.Do(t, context.Background()).Status(t, action.AssertEquals(http.StatusOK))
		})
	}
}

func TestRequest_Body(t *testing.T) {
	for _, item := range []struct {
		name     string
		body     any
		expected string
	}{
		{
			name:     "test string",
			body:     `{"action":"action"}`,
			expected: `{"action":"action"}`,
		},
		{
			name:     "test string",
			body:     []byte(`{"action":"action"}`),
			expected: `{"action":"action"}`,
		},
		{
			name:     "test string",
			body:     strings.NewReader(`{"action":"action"}`),
			expected: `{"action":"action"}`,
		},

		{
			name: "test struct",
			body: struct {
				Value string `json:"action"`
			}{
				Value: "action",
			},
			expected: `{"action":"action"}`,
		},
	} {
		t.Run(item.name, func(t *testing.T) {
			r := newRequest(t, http.MethodGet, "http://localhost:8080", &Deps{Handler: noopHandler{}})
			r.Body(t, item.body)

			assert.JSONEq(t, item.expected, string(must(io.ReadAll(r.req.Body))))
		})
	}

}
