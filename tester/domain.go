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
	"github.com/phonkee/jayson/tester/resolver"
	"github.com/stretchr/testify/require"
	"net/http"
)

// APIClient is the interface for testing the rest APIClient
type APIClient interface {
	// Delete does a DELETE response to the APIClient
	Delete(t require.TestingT, path string) APIRequest

	// Get does a GET response to the APIClient
	Get(t require.TestingT, path string) APIRequest

	// Post does a POST response to the APIClient
	Post(t require.TestingT, path string) APIRequest

	// Put does a PUT response to the APIClient
	Put(t require.TestingT, path string) APIRequest

	// Request does a response to the APIClient
	Request(t require.TestingT, method string, path string) APIRequest

	// ReverseURL creates a path by given url name and url arguments
	ReverseURL(t require.TestingT, name string, extra ...resolver.Extra) string
}

// APIRequest is the interface for testing the rest APIClient response
type APIRequest interface {
	// Body sets the body of the request
	Body(t require.TestingT, body any) APIRequest
	// Do perform the request and returns the response
	Do(t require.TestingT, ctx context.Context) APIResponse
	// Header sets the header of the request
	Header(t require.TestingT, key, value string) APIRequest
	// Query sets the query of the request
	Query(t require.TestingT, key, value string) APIRequest
}

// APIResponse is the interface for testing the rest APIClient response
type APIResponse interface {
	// AssertHeaderValue asserts that response header has given value for given key
	AssertHeaderValue(t require.TestingT, key, value string) APIResponse

	// AssertJsonEquals asserts that response body is equal to given json
	AssertJsonEquals(t require.TestingT, what any) APIResponse

	// AssertJsonPath asserts that given json path conforms to given value (special operations are supported)
	// This method will unmarshal the body into the same type as the value
	AssertJsonPath(t require.TestingT, path string, what any) APIResponse

	// AssertStatus asserts that response status is equal to given status
	AssertStatus(t require.TestingT, status int) APIResponse

	// Unmarshal json unmarshalls whole response body into given value
	Unmarshal(t require.TestingT, v any) APIResponse
}

const (
	// ContentTypeHeader is the content type header
	ContentTypeHeader = "Content-Type"
	// ContentTypeJSON is the content type for json
	ContentTypeJSON = "application/json"
)

//go:generate mockery --name=RoundTripper --filename=round_tripper.go
type RoundTripper interface {
	http.RoundTripper
}
