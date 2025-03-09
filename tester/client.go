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
	"github.com/phonkee/jayson/tester/resolver"
	"github.com/stretchr/testify/require"
	"net/http"
)

// WithAPI runs the given function with a new APIClient instance
func WithAPI(t require.TestingT, deps *Deps, fn func(APIClient)) {
	// validate deps first
	deps.Validate(t)

	// call closure with an api client
	fn(newClient(t, deps))
}

// newClient creates a new *client instance
func newClient(t require.TestingT, deps *Deps) *client {
	return &client{
		deps: deps,
	}
}

// client is the implementation of rest APIClient testing helper
type client struct {
	deps *Deps
}

// Delete does a DELETE response to the APIClient
func (c *client) Delete(t require.TestingT, path string) APIRequest {
	return c.Request(t, http.MethodDelete, path)
}

// Get does a GET response to the APIClient
func (c *client) Get(t require.TestingT, path string) APIRequest {
	return c.Request(t, http.MethodGet, path)
}

// Post does a POST response to the APIClient
func (c *client) Post(t require.TestingT, path string) APIRequest {
	return c.Request(t, http.MethodPost, path)
}

// Put does a PUT response to the APIClient
func (c *client) Put(t require.TestingT, path string) APIRequest {
	return c.Request(t, http.MethodPut, path)
}

// Request does a response to the APIClient
func (c *client) Request(t require.TestingT, method string, path string) APIRequest {
	return newRequest(method, path, c.deps)
}

// ReverseURL creates a path by given url name and resolver extra
func (c *client) ReverseURL(t require.TestingT, name string, extra ...resolver.Extra) string {
	return c.deps.ReverseURL(t, name, extra...)
}

// ReverseArgs adds arguments key value pairs to resolver.Extra for ReverseURL
func (c *client) ReverseArgs(t require.TestingT, kv ...string) resolver.Extra {
	return c.deps.ReverseArgs(t, kv...)
}

// ReverseQuery adds query key value pairs to resolver.Extra for ReverseURL
func (c *client) ReverseQuery(t require.TestingT, kv ...string) resolver.Extra {
	return c.deps.ReverseQuery(t, kv...)
}
