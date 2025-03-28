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
	"github.com/phonkee/jayson/tester/resolver"
	"github.com/stretchr/testify/require"
	"io"
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

	// ReverseArgs creates a resolver.Extra from given key value pairs
	ReverseArgs(t require.TestingT, kv ...string) resolver.Extra

	// ReverseQuery creates a resolver.Extra from given key value pairs
	ReverseQuery(t require.TestingT, kv ...string) resolver.Extra
}

// APIRequest is the interface for testing the rest APIClient response
type APIRequest interface {
	// Body sets the body of the request
	Body(t require.TestingT, body any) APIRequest

	// Do perform the request and returns the response
	Do(t require.TestingT, ctx context.Context) APIResponse

	// Header sets the header of the request
	Header(t require.TestingT, key, value string) APIRequest

	// Print the response to stdout
	Print(writer ...io.Writer) APIRequest

	// Query sets the query of the request
	Query(t require.TestingT, key, value string) APIRequest
}

// APIResponse is the interface for testing the rest APIClient response
type APIResponse interface {

	// Print the response to stdout
	Print(writer ...io.Writer) APIResponse

	// Header asserts that response header action is equal to given action
	Header(t require.TestingT, key string, action action.Action) APIResponse

	// Json calls the action on the json response path
	Json(t require.TestingT, path string, action action.Action) APIResponse

	// Status asserts that response status is equal to given action
	Status(t require.TestingT, action action.Action) APIResponse
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

const (
	// RootPath is the root path for Response Json call
	RootPath string = ""
)
