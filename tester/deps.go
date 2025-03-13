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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/url"
	"strings"
)

// Deps is the dependencies for the APIClient
// Resolve provides ReverseURL method, but is only optional
// One of Handler or Address is required
type Deps struct {
	// Handler is the http.Handler
	Handler http.Handler
	// Addr is the address of the server (host:port)
	Address string
	// Client sets custom http client
	Client *http.Client
	// Resolver is the URL resolver
	Resolver resolver.Resolver
}

// ReverseURL resolves URL by name
func (d *Deps) ReverseURL(t require.TestingT, name string, extra ...resolver.Extra) string {
	assert.NotNil(t, d.Resolver, "Resolver is required")
	return d.Resolver.ReverseURL(t, name, extra...)
}

// ReverseArgs adds arguments key action pairs to resolver.Extra for ReverseURL
func (d *Deps) ReverseArgs(t require.TestingT, kv ...string) resolver.Extra {
	assert.NotNil(t, d.Resolver, "Resolver is required")
	return resolver.Arguments(t, kv...)
}

// ReverseQuery adds query key action pairs to resolver.Extra for ReverseURL
func (d *Deps) ReverseQuery(t require.TestingT, kv ...string) resolver.Extra {
	assert.NotNil(t, d.Resolver, "Resolver is required")
	return resolver.Query(t, kv...)
}

// Validate deps
func (d *Deps) Validate(t require.TestingT) {
	// clean address
	if d.Address = strings.TrimSpace(d.Address); d.Address != "" {
		// parse url here
		if parsed, err := url.Parse(d.Address); err == nil && parsed.Host != "" {
			d.Address = parsed.Host
		}
	}
	// check if exampleHandler or Address is provided
	require.Falsef(t, d.Handler == nil && d.Address == "", "Handler or Address is required")
}
