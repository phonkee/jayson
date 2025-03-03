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
	"testing"
)

// APIClient is the interface for testing the rest APIClient
type APIClient interface {
	// Delete does a DELETE response to the APIClient
	Delete(t *testing.T, path string) APIRequest

	// Get does a GET response to the APIClient
	Get(t *testing.T, path string) APIRequest

	// Post does a POST response to the APIClient
	Post(t *testing.T, path string) APIRequest

	// Put does a PUT response to the APIClient
	Put(t *testing.T, path string) APIRequest

	// Request does a response to the APIClient
	Request(t *testing.T, method string, path string) APIRequest

	// ReverseURL creates a path by given url name and url arguments
	ReverseURL(t *testing.T, name string, vars ...string) string
}

// APIRequest is the interface for testing the rest APIClient response
type APIRequest interface {
	// Body sets the body of the request
	Body(t *testing.T, body any) APIRequest
	// Do perform the request and returns the response
	Do(t *testing.T, ctx context.Context) APIResponse
	// Header sets the header of the request
	Header(t *testing.T, key, value string) APIRequest
	// Query sets the query of the request
	Query(t *testing.T, key, value string) APIRequest
}

// APIResponse is the interface for testing the rest APIClient response
type APIResponse interface {
	// AssertJsonEquals asserts that response body is equal to given json
	AssertJsonEquals(t *testing.T, what any) APIResponse

	// AssertJsonKeyEquals asserts that response body key is equal to given value
	AssertJsonKeyEquals(t *testing.T, key string, what any) APIResponse

	// AssertStatus asserts that response status is equal to given status
	AssertStatus(t *testing.T, status int) APIResponse

	// Unmarshal json unmarshalls whole response body into given value
	Unmarshal(t *testing.T, v any) APIResponse
}
