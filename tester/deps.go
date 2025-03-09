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
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/url"
	"strings"
)

// Deps is the dependencies for the APIClient
// Router is optional, if not provided, ReverseURL will not work
// One of Handler or Address is required
type Deps struct {
	// Router - currently required
	Router *mux.Router
	// Handler is the http.Handler
	Handler http.Handler
	// Addr is the address of the server
	Address string
	// Client sets custom http client
	Client *http.Client
}

// Validate deps
func (d *Deps) Validate(t require.TestingT) {
	// clean address
	if d.Address = strings.TrimSpace(d.Address); d.Address != "" {
		// parse url here
		if parsed, err := url.Parse(d.Address); err == nil {
			d.Address = parsed.Host
		}
	}
	// check if exampleHandler or Address is provided
	require.Falsef(t, d.Handler == nil && d.Address == "", "exampleHandler or Address is required")
}
